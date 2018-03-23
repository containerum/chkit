package update

import (
	"io"

	"fmt"

	"runtime"

	"net/http"

	"github.com/cheggaaa/pb"
	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
	"gopkg.in/urfave/cli.v2"
)

type LatestChecker interface {
	LatestVersion() (string, error)
}

type LatestDownloader interface {
	LatestDownload() (io.ReadCloser, error)
}

type LatestCheckerDownloader interface {
	LatestChecker
	LatestDownloader
}

const (
	ErrUpdateCheck    = chkitErrors.Err("unable to check latest version")
	ErrUpdateDownload = chkitErrors.Err("unable to download latest version")
)

func DownloadFileName(version string) string {
	extension := "tar.gz"
	if runtime.GOOS == "windows" {
		extension = "zip"
	}
	return fmt.Sprintf("chkit_%s_%s_v%s.%s", runtime.GOOS, runtime.GOARCH, version, extension)
}

type GithubLatestCheckerDownloader struct {
	client      *resty.Client
	log         *logrus.Logger
	ctx         *cli.Context
	downloadUrl string
}

func NewGithubLatestCheckerDownloader(ctx *cli.Context, owner, repo string) *GithubLatestCheckerDownloader {
	return &GithubLatestCheckerDownloader{
		ctx: ctx,
		log: util.GetLog(ctx),
		client: resty.New().
			SetHostURL(fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo)).
			SetDebug(true).
			SetLogger(util.GetLog(ctx).WriterLevel(logrus.DebugLevel)),
		downloadUrl: fmt.Sprintf("https://github.com/%s/%s/releases/download", owner, repo),
	}
}

func (gh *GithubLatestCheckerDownloader) LatestVersion() (string, error) {
	gh.log.Debug("get latest version from github")

	var latestVersionResp struct {
		LatestVersion string `json:"tag_name"`
	}

	_, err := gh.client.R().SetResult(&latestVersionResp).Get("/latest")
	if err != nil {
		return "0.0.1-alpha", chkitErrors.Wrap(ErrUpdateCheck, err)
	}

	return latestVersionResp.LatestVersion, nil
}

func (gh *GithubLatestCheckerDownloader) LatestDownload() (io.ReadCloser, error) {
	gh.log.Debug("download update")

	latestVersion, err := gh.LatestVersion()
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(fmt.Sprintf("%s/%s/%s", gh.downloadUrl, latestVersion, DownloadFileName(latestVersion)))
	if err != nil {
		return nil, chkitErrors.Wrap(ErrUpdateDownload, err)
	}

	bar := pb.New64(resp.ContentLength)

	return bar.NewProxyReader(resp.Body), nil
}
