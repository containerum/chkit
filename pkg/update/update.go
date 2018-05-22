package update

import (
	"bufio"
	"crypto"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/blang/semver"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/context"
	"github.com/inconshreveable/go-update"
	"github.com/sirupsen/logrus"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"golang.org/x/crypto/ssh/terminal"
)

var PublicKeyB64 = "cHVibGljIGtleQo="

const (
	ErrUpdateApply = chkitErrors.Err("update apply failed")
)

func verifiedUpdate(upd *Package) error {
	checksum, err := ioutil.ReadAll(upd.Hash)
	if err != nil {
		return chkitErrors.Wrap(ErrUpdateApply, err)
	}

	signature, err := ioutil.ReadAll(upd.Signature)
	if err != nil {
		return chkitErrors.Wrap(ErrUpdateApply, err)
	}

	opts := update.Options{
		Checksum:  checksum,
		Signature: signature,
		Verifier:  update.NewECDSAVerifier(),
		Hash:      crypto.SHA256,
	}
	publicKey, err := base64.StdEncoding.DecodeString(PublicKeyB64)
	if err != nil {
		return chkitErrors.Wrap(ErrUpdateApply, err)
	}

	err = opts.SetPublicKeyPEM(publicKey)
	if err != nil {
		return chkitErrors.Wrap(ErrUpdateApply, err)
	}

	err = update.Apply(upd.Binary, opts)
	if err != nil {
		return chkitErrors.Wrap(ErrUpdateApply, err)
	}

	return nil
}

func AskForUpdate(ctx *context.Context, latestVersion semver.Version) (bool, error) {
	// check if we have terminal supports escape sequences
	var colorStart, colorEnd string
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		colorStart = "\x1b[31;1m"
		colorEnd = "\x1b[0m"
	}
	fmt.Printf("%sYou are using version %s, however version %s is available%s\n",
		colorStart, ctx.Version, latestVersion, colorEnd)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	for {
		fmt.Println("Do you want to update [Y/n]: ")
		for !scanner.Scan() {
			break
		}
		if err := scanner.Err(); err != nil {
			logrus.WithError(err).Error("scan failed")
			return false, err
		}
		switch strings.ToLower(scanner.Text()) {
		case "", "y":
			return true, nil
		case "n":
			return false, nil
		default:
			continue
		}
	}
}

func Update(downloader LatestCheckerDownloader, restartAfter bool) error {
	archive, size, err := downloader.LatestDownload()
	if err != nil {
		return err
	}
	defer archive.Close()

	p := mpb.New()
	bar := p.AddBar(size, mpb.PrependDecorators(
		decor.Counters("%.1f / %.1f", 1, 3, 0),
	), mpb.AppendDecorators(
		decor.Percentage(3, 0),
		decor.StaticName(" ETA:", 3, 0),
		decor.ETA(3, 0),
	))
	archive = bar.ProxyReader(archive)

	pkg, err := unpack(archive)
	if err != nil {
		return err
	}
	p.Wait()

	err = verifiedUpdate(pkg)
	if err != nil {
		return err
	}

	if restartAfter {
		gracefulRestart()
	}

	return nil
}
