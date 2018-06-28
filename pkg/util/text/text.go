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

func Crop(txt string, width int) string {
	split := strings.Split(txt, "\n")
	if len(split) > 1 {
		txt = split[0] + "..."
	}
	txtWidth := Width(txt)
	if txtWidth <= width {
		return txt
	}
	if width <= 3 {
		return string([]rune(txt)[:width])
	}
	return string([]rune(txt)[:width-3]) + "..."
}

func Indent(text string, indent uint) string {
	lines := strings.Split(text, "\n")
	ind := strings.Repeat(" ", int(indent))
	for i, line := range lines {
		lines[i] = ind + line
	}
	return strings.Join(lines, "\n")
}
