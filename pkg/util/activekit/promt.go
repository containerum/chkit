package activekit

import "fmt"

func Promt(promt string) string {
	fmt.Print(promt)
	return Input()
}
