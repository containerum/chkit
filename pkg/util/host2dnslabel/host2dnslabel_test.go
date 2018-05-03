package host2dnslabel

import (
	"testing"

	"github.com/containerum/chkit/pkg/util/validation"
)

func TestHost2DNSLabel(test *testing.T) {
	var hosts = []string{
		"google.com",
		"123.com",
		"asndl0-😏ю.loc",
		"as=-0e20 -doqd- 3-- -s.saalc=asd.cpks",
		"asdasd d-ds d----- --.net",
	}
	for _, host := range hosts {
		DNSlabel := Host2DNSLabel(host)
		test.Logf("%q -> %q", host, DNSlabel)
		if err := validation.DNSLabel(DNSlabel); err != nil {
			test.Fatal(err)
		}
	}
}
