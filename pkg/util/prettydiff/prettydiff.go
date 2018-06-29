package prettydiff

import (
	"io"
	"strings"
)

func Fprint(wr io.Writer, diff string) error {
	for _, line := range strings.Split(diff, "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "-"):
			line = "\x1b[31;1m" + line + "\x1b[0m"
		case strings.HasPrefix(trimmed, "+"):
			line = "\x1b[32m" + line + "\x1b[0m"
		case strings.HasPrefix(trimmed, "@"):
			line = "\x1b[34m" + line + "\x1b[0m"
		}
		if _, err := wr.Write([]byte(line + "\n")); err != nil {
			return err
		}
	}
	return nil
}
