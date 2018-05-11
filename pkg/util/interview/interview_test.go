package interview

import (
	"bytes"
	"crypto/rand"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
	"testing"
	"time"
)

func TestView(test *testing.T) {
	var currCase *testCase
	defer func() {
		if rec := recover(); rec != nil {
			test.Fatalf("PANIC in test %q:\n%v\n%s\n", currCase, rec, debug.Stack())
		}
	}()
	var randomBytes = func() []byte {
		var buf = make([]byte, 16)
		rand.Read(buf)
		return buf
	}()
	var err = errors.New("error message")
	var testStruct = struct {
		Foo string
	}{
		Foo: "bar",
	}
	now := time.Now()
	var testData = []testCase{
		{In: nil, Expected: strPtr("nil")},
		{In: (*error)(nil), Expected: strPtr("<nil>")},
		{In: error(nil), Expected: strPtr("nil")},
		{In: "string", Expected: strPtr("string")},
		{In: int64(42), Expected: strPtr("42")},
		{In: []byte("hecking shiba"), Expected: strPtr("hecking shiba")},
		{In: randomBytes, Expected: strPtr(fmt.Sprintf("%x", randomBytes))},
		{In: err, Expected: strPtr(err.Error())},
		{In: bytes.NewBufferString("buffer"), Expected: strPtr("buffer")},
		{In: testStruct, Expected: strPtr(fmt.Sprintf("%#v", testStruct))},
		{In: 4.2, Expected: strPtr("4.2")},
		{In: jsonMarshalTest{now}, Expected: func() *string {
			data, err := now.MarshalJSON()
			if err != nil {
				panic(err)
			}
			return strPtr(string(data))
		}()},
		{In: textMarshalTest{now}, Expected: func() *string {
			data, err := now.MarshalText()
			if err != nil {
				panic(err)
			}
			return strPtr(string(data))
		}()},
		{In: &testStruct, Expected: strPtr(fmt.Sprintf("%#v", testStruct))},
		{In: jsonMarshalTest{marshalErr{}}, Expected: strPtr(fallbackView(jsonMarshalTest{marshalErr{}}))},
		{In: textMarshalTest{marshalErr{}}, Expected: strPtr(fallbackView(textMarshalTest{marshalErr{}}))},
	}
	for i, testCase := range testData {
		currCase = &testCase
		currCase.ID = i
		var out = View(testCase.In)
		if testCase.Expected != nil && out != *testCase.Expected {
			test.Errorf("TEST CASE %d expected %q, got %#v", i, *testCase.Expected, out)
		}
	}
}

func strPtr(str string) *string {
	return &str
}

type testCase struct {
	ID       int
	Name     string
	In       interface{}
	Expected *string
}

func (t *testCase) String() string {
	var str string
	str += fmt.Sprintf("%d ", t.ID)
	if t.Name != "" {
		str += t.Name + " -> "
	}
	str += "In: " + fmt.Sprint(t.In)
	if t.Expected != nil {
		str += ", Expected: " + *t.Expected
	}
	return str
}

type jsonMarshalTest struct {
	json.Marshaler
}

type textMarshalTest struct {
	encoding.TextMarshaler
}

type marshalErr struct{}

func (marshalErr) MarshalJSON() ([]byte, error) {
	return nil, errors.New("errorial json marshaler")
}

func (marshalErr) MarshalText() ([]byte, error) {
	return nil, errors.New("errorial text marshaler")
}
