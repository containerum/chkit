package helpers

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

type MappedStruct map[string]interface{}

func tagForField(fieldType reflect.StructField) string {
	tag, ok := fieldType.Tag.Lookup("mapconv")
	if tag == "-" {
		return ""
	}
	if !ok {
		tag = fieldType.Name
	}
	return tag
}

func StructToMap(s interface{}) (ret MappedStruct) {
	ret = make(MappedStruct)
	structRefl := reflect.Indirect(reflect.ValueOf(s))
	structReflType := structRefl.Type()
	for i := 0; i < structReflType.NumField(); i++ {
		field := structRefl.Field(i)
		fieldType := structReflType.Field(i)
		tag := tagForField(fieldType)
		if tag == "" {
			continue
		}
		if fieldType.Type.Kind() == reflect.Struct {
			ret[tag] = StructToMap(field.Interface())
		} else {
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			enc.EncodeValue(field)
			ret[tag] = buf.Bytes()
		}
	}
	return
}

func FillStruct(s interface{}, data MappedStruct) error {
	structRefl := reflect.Indirect(reflect.ValueOf(s))
	structReflType := structRefl.Type()
	for i := 0; i < structReflType.NumField(); i++ {
		field := structRefl.Field(i)
		fieldType := structReflType.Field(i)
		tag := tagForField(fieldType)
		if tag == "" {
			continue
		}
		fieldData, hasFd := data[tag]
		if !hasFd {
			continue
		}
		if fieldType.Type.Kind() == reflect.Struct {
			err := FillStruct(field.Interface(), fieldData.(MappedStruct))
			if err != nil {
				return err
			}
		} else {
			buf := bytes.NewBuffer(fieldData.([]byte))
			dec := gob.NewDecoder(buf)
			dec.DecodeValue(field)
		}
	}
	return nil
}
