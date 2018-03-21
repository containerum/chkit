package delog

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	gopaths = []string{}
)

func stripGopath(str string) string {
	str = sanitizePath(str)
	for _, gopath := range gopaths {
		if strings.HasPrefix(str, gopath) {
			return str[len(gopath):]
		}
	}
	return str
}

func sanitizePath(str string) string {
	str = filepath.ToSlash(str)
	probablyDiskAndGopath := strings.SplitN(str, ":", 2)
	if len(probablyDiskAndGopath) > 1 {
		str = strings.ToLower(probablyDiskAndGopath[0]) + ":" + probablyDiskAndGopath[1]
	}
	return str
}

func init() {
	goroot := runtime.GOROOT()
	gopaths = []string{path.Join(filepath.ToSlash(goroot), "src") + "/"}
	for _, gopath := range filepath.SplitList(os.Getenv("GOPATH")) {
		if gopath != "" {
			gopath = sanitizePath(gopath)
			gopaths = append(gopaths, path.Join(gopath, "src")+"/")
		}
	}
}
