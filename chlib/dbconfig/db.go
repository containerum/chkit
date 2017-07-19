package dbconfig

import (
	"chkit-v2/helpers"
	"fmt"
	"github.com/boltdb/bolt"
	"os"
)

var initializers map[string]helpers.MappedStruct = make(map[string]helpers.MappedStruct)

type ConfigDB struct {
	db *bolt.DB
}

func (d *ConfigDB) initialize(bucket *bolt.Bucket, initializer helpers.MappedStruct) error {
	for key, value := range initializer {
		switch value.(type) {
		case []byte: // simple value, just put to bucket
			err := bucket.Put([]byte(key), value.([]byte))
			if err != nil {
				return err
			}
		case helpers.MappedStruct: // structure data, recursive put to bucket
			newBucket, err := bucket.CreateBucketIfNotExists([]byte(key))
			if err != nil {
				return err
			}
			return d.initialize(newBucket, value.(helpers.MappedStruct))
		}
	}
	return nil
}

func OpenOrCreate(filePath string) (db *ConfigDB, err error) {
	db = new(ConfigDB)
	db.db, err = bolt.Open(filePath, os.ModePerm, nil)
	if err != nil {
		return
	}
	// for all top-level buckets
	err = db.db.Update(func(tx *bolt.Tx) error {
		// for all top-level buckets
		for name, initializer := range initializers {
			bucket := tx.Bucket([]byte(name))
			if bucket == nil {
				bucket, err := tx.CreateBucketIfNotExists([]byte(name))
				if err != nil {
					return err
				}
				err = db.initialize(bucket, initializer)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	return
}

func (d *ConfigDB) readRecursive(bucket *bolt.Bucket) helpers.MappedStruct {
	data := make(helpers.MappedStruct)
	bucket.ForEach(func(k, v []byte) error {
		if v != nil { // simple key-value
			data[string(k)] = v
		} else { // bucket - has underlying structure
			newBucket := bucket.Bucket(k)
			data[string(k)] = d.readRecursive(newBucket)
		}
		return nil
	})
	return data
}

func (d *ConfigDB) readTransactional(bucketName string) (m helpers.MappedStruct, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("transactional read: no such bucket %s", bucketName)
		}
		m = d.readRecursive(bucket)
		return nil
	})
	return
}

func (d *ConfigDB) writeRecursive(m helpers.MappedStruct, bucket *bolt.Bucket) (err error) {
	for key, value := range m {
		switch value.(type) {
		case []byte:
			err = bucket.Put([]byte(key), value.([]byte))
		case helpers.MappedStruct:
			newBucket, err := bucket.CreateBucketIfNotExists([]byte(key))
			if err != nil {
				return err
			}
			err = d.writeRecursive(value.(helpers.MappedStruct), newBucket)
		}
		if err != nil {
			return fmt.Errorf("value for key %s push: %s", key, err)
		}
	}
	return
}

func (d *ConfigDB) writeTransactional(m helpers.MappedStruct, bucketName string) (err error) {
	return d.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("transactional read: %s", err)
		}
		return d.writeRecursive(m, bucket)
	})
}

func (d *ConfigDB) Close() {
	d.db.Close()
}
