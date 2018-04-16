package text

import (
	"strings"
	"unicode/utf8"
)

func Width(text string) int {
	width := 0
	for _, line := range strings.Split(text, "\n") {
		l := utf8.RuneCountInString(line)
		if l > width {
			width = l
		}
	}
	return width
}
