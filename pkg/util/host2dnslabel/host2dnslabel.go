package host2dnslabel

import (
	"bytes"
	"strings"
	"unicode"

	"github.com/containerum/chkit/pkg/util/namegen"
)

func Host2DNSLabel(hostname string) string {
	hostname = strings.ToLower(hostname)
	var label = bytes.NewBuffer(make([]byte, 0, len(hostname)))
	var dash = false
	for _, ch := range hostname {
		switch {
		case unicode.IsPunct(ch):
			ch = '-'
			if !dash {
				label.WriteRune(ch)
			}
			dash = true
		case (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9'):
			label.WriteRune(ch)
			dash = false
		default:
			continue
		}
	}
	if label.Len() == 0 {
		return namegen.Aster() + "-" + namegen.Color() + "-" + namegen.Physicist()
	}
	return label.String()
}
