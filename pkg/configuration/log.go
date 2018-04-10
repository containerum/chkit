package configuration

import "time"

func LogFileName() string {
	return time.Now().Format("2006-Jan-2") + ".log"
}
