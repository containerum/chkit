package activekit

import (
	"fmt"
	"strings"
)

func YesNo(promt string) bool {
	fmt.Print(promt + " [Y/N]: ")
	answer := strings.ToLower(strings.TrimSpace(Input()))
	return answer == "y"
}
