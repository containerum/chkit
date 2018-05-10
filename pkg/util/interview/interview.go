package interview

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"unicode/utf8"
)

func View(v interface{}) (str string) {
	switch value := v.(type) {
	case nil:
		return "nil"
	case fmt.Stringer:
		return value.String()
	case error:
		return value.Error()
	case json.Marshaler:
		data, err := value.MarshalJSON()
		if err != nil {
			return fallbackView(value)
		}
		return string(data)
	case encoding.TextMarshaler:
		txt, err := value.MarshalText()
		if err != nil {
			return fallbackView(value)
		}
		return string(txt)
	}
	var value = reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Interface, reflect.Ptr:
		if value.IsNil() {
			return "<nil>"
		}
		value = value.Elem()
	}
	switch value := value.Interface().(type) {
	case string:
		return value
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", value)
	case float32, float64, complex64, complex128:
		return fmt.Sprintf("%g", value)
	case []byte:
		if utf8.ValidString(string(value)) {
			return string(value)
		}
		return fmt.Sprintf("%x", value)

	default:
		return fallbackView(value)
	}
}

func fallbackView(v interface{}) string {
	return fmt.Sprintf("%#v", v)
}
