package ferr

import (
	"fmt"
	"os"
)

func Println(args ...interface{}) (int, error) {
	return fmt.Fprint(os.Stderr, args...)
}

func Printf(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(os.Stderr, format, args...)
}

func Print(args ...interface{}) (int, error) {
	return fmt.Fprint(os.Stderr, args...)
}
