package activekit

import (
	"fmt"
	"strings"
)

func YesNo(promt string, args ...interface{}) bool {
	fmt.Printf("%s [Y/N]: ", fmt.Sprintf(promt, args...))
	answer := strings.ToLower(strings.TrimSpace(Input()))
	return answer == "y"
}
