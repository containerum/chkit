package chlib

import (
	"fmt"
	"os/user"
	"path"

	"github.com/boltdb/bolt"
	"github.com/kfeofantov/chkit-v2/helpers"
)

const configDir = ".containerum"
const configFile = "config.db"

var configPath string

var configDb *bolt.DB

var initializers map[string]helpers.MappedStruct = make(map[string]helpers.MappedStruct)

func init() {
	currentUser, err := user.Current()
	if err != nil {
		panic(fmt.Errorf("get current user: %s", err))
	}

	configPath = path.Join(currentUser.HomeDir, configDir)
	configDb, err = bolt.Open(path.Join(configPath, configFile), 0600, nil)
	if err != nil {
		panic(fmt.Errorf("config db open: %s", err))
	}
	if err := runInitializers(); err != nil {
		panic(err)
	}
}

func recursiveInitialize(startBucket *bolt.Bucket, initializer helpers.MappedStruct) error {
	for key, value := range initializer {
		switch value.(type) {
		case []byte:
			err := startBucket.Put([]byte(key), value.([]byte))
			if err != nil {
				return err
			}
		case helpers.MappedStruct:
			newBucket, err := startBucket.CreateBucketIfNotExists([]byte(key))
			if err != nil {
				return fmt.Errorf("bucket create: %s", err)
			}
			err = recursiveInitialize(newBucket, value.(helpers.MappedStruct))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func runInitializers() error {
	return configDb.Update(func(tx *bolt.Tx) error {
		for bucketName, initializer := range initializers {
			bucket := tx.Bucket([]byte(bucketName))
			if bucket == nil {
				bucket, err := tx.CreateBucket([]byte(bucketName))
				if err != nil {
					return fmt.Errorf("bucket create: %s", err)
				}
				recursiveInitialize(bucket, initializer)
			}
		}
		return nil
	})
}

func recursivePush(startBucket *bolt.Bucket, data helpers.MappedStruct) error {
	for key, value := range data {
		var err error
		switch value.(type) {
		case []byte:
			err = startBucket.Put([]byte(key), value.([]byte))
		case helpers.MappedStruct:
			newBucket, err := startBucket.CreateBucket([]byte(key))
			if err != nil {
				return err
			}
			err = recursivePush(newBucket, value.(helpers.MappedStruct))
		}
		if err != nil {
			return fmt.Errorf("value for key %s push: %s", key, err)
		}
	}
	return nil
}

func pushToBucket(bucketName string, data helpers.MappedStruct) error {
	return configDb.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		return recursivePush(bucket, data)
	})
}

func recursiveRead(startBucket *bolt.Bucket) (data helpers.MappedStruct, err error) {
	data = make(helpers.MappedStruct)
	err = startBucket.ForEach(func(k []byte, v []byte) error {
		if v != nil {
			data[string(k)] = v
		} else {
			// it`s a bucket
			newBucket := startBucket.Bucket(k)
			var err error
			data[string(k)], err = recursiveRead(newBucket)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return
}

func readFromBucket(bucketName string) (data helpers.MappedStruct, err error) {
	err = configDb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s not exist", bucketName)
		}
		data, err = recursiveRead(bucket)
		return err
	})
	return
}

func Close() {
	configDb.Close()
}
