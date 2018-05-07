package pairs

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMap(test *testing.T) {
	const delim = ":"
	var origin = map[string]string{
		"HOME": "/home/dir",
		"PATH": "/plan9/buzz",
		"STU":  "/pearl/jam",
		"X":    "212321312",
	}

	var kvStr = func() string {
		buf := &bytes.Buffer{}
		for k, v := range origin {
			fmt.Fprintf(buf, "%s%s%q ", k, delim, v)
		}
		return buf.String()
	}()
	pairs, err := ParseMap(kvStr, delim)
	if err != nil {
		test.Fatal(err)
	}
	if !assert.Equal(test, origin, pairs) {
		test.Fail()
	}
}
