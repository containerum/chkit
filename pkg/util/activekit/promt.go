package activekit

import "fmt"

func Promt(promt string, vars ...interface{}) string {
	fmt.Printf(promt, vars...)
	return Input()
}
