package angel

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"runtime/debug"

	"github.com/containerum/chkit/pkg/configdir"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
)

func reportPreambula(version string) string {
	return fmt.Sprintf("[REPORT]\n"+
		"chkit fatal report\n"+
		"version: %s\n"+
		"os: %s %s", version, runtime.GOOS, runtime.GOARCH)
}

func Angel(ctx *context.Context, sin interface{}) {
	report := reportPreambula(ctx.Version)
	switch recoverData := sin.(type) {
	case nil:
		return
	case error:
		report += fmt.Sprintf("[FATAL] %v", recoverData)
	default:
		report += fmt.Sprintf("[FATAL] %v\n%s", recoverData, string(debug.Stack()))
	}
	logDir := configdir.LogDir()
	logFileName := configuration.LogFileName()
	reportFile := path.Join(logDir, "report.txt")

	err := ioutil.WriteFile(reportFile, []byte(report), os.ModePerm)
	if err != nil {
		fmt.Printf("[FATAL] something completely wrong.\n")
		fmt.Printf("Please, send %q and %q files from %q to support@exonlab.omnidesk.ru\n",
			logFileName, "report.txt", logDir)
		return
	}
	/*
		logTail, err := readLogTail(logFilePath)
		if err != nil {
			fmt.Printf("[FATAL] something completely wrong.\n")
			fmt.Printf("Please, send report and log file from %q to support@exonlab.omnidesk.ru", configDir)
			return
		}

		report = report + logTail
		if err := openSupportPageWithReport(report); err != nil {
	*/
	fmt.Printf("Fatal error: %v\n", sin)
	fmt.Printf("Please, send %q and %q files from %q to support@exonlab.omnidesk.ru\n",
		logFileName, "report.txt", logDir)
	//	}
}
