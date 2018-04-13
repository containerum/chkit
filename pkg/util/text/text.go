package text

import "strings"

func Width(text string) int {
	width := 0
	for _, line := range strings.Split(text, "\n") {
		if len(line) > width {
			width = len(line)
		}
	}
	return width
}
