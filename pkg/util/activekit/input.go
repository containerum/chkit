package activekit

import (
	"bufio"
	"os"
)

func Input() string {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		return scanner.Text()
	}
	return ""
}
