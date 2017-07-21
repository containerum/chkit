package helpers

import (
	"os"
	"path/filepath"
	"time"

	"github.com/kardianos/osext"
)

var CurrentClientVersion string

func GetProgramBuildTime() time.Time {
	fname, _ := osext.Executable()
	dir, err := filepath.Abs(filepath.Dir(fname))
	if err != nil {
		return time.Now()
	}
	fi, err := os.Lstat(filepath.Join(dir, filepath.Base(fname)))
	if err != nil {
		return time.Now()
	}
	return fi.ModTime()
}
