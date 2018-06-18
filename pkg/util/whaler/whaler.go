package whaler

import (
	"fmt"
	"regexp"
)

var (
	imageRe = regexp.MustCompile("(?:.+/)?([^:]+)(?::.+)?")
)

func ExtractImageName(name string) (string, error) {
	var tokens = imageRe.FindAllStringSubmatch(name, -1)
	if len(tokens) > 0 && len(tokens[0]) > 1 {
		return tokens[0][1], nil
	} else {
		return "", fmt.Errorf("invalid container name %q, must match with %q", name, imageRe.String())
	}
}
