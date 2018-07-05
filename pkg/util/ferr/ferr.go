package ferr

import (
	"fmt"
	"os"
	"sync"
)

// global IO lock
var giol sync.Mutex

func Println(args ...interface{}) (int, error) {
	giol.Lock()
	defer giol.Unlock()
	return fmt.Fprint(os.Stderr, args...)
}

func Printf(format string, args ...interface{}) (int, error) {
	giol.Lock()
	defer giol.Unlock()
	return fmt.Fprintf(os.Stderr, format, args...)
}

func Print(args ...interface{}) (int, error) {
	giol.Lock()
	defer giol.Unlock()
	return fmt.Fprint(os.Stderr, args...)
}
