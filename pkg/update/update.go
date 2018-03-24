package update

import (
	"io/ioutil"

	"crypto"

	"fmt"

	"os"

	"bufio"

	"strings"

	"encoding/base64"

	"github.com/blang/semver"
	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/inconshreveable/go-update"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/urfave/cli.v2"
)

var PublicKeyB64 = "cHVibGljIGtleQo="

const (
	ErrUpdateApply  = chkitErrors.Err("update apply failed")
	ErrVersionParse = chkitErrors.Err("version parse failed")
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

func Update(ctx *cli.Context, downloader LatestCheckerDownloader, restartAfter bool) error {
	latestVersion, err := downloader.LatestVersion()
	if err != nil {
		return err
	}

	latestVersionSemver, err := semver.ParseTolerant(latestVersion)
	if err != nil {
		return chkitErrors.Wrap(ErrVersionParse, err)
	}

	if latestVersionSemver.LE(util.GetVersion(ctx)) {
		return nil
	}

	// check if we have terminal supports escape sequences
	var colorStart, colorEnd string
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		colorStart = "\x1b[31;1m"
		colorEnd = "\x1b[0m"
	}
	fmt.Printf("%sYou are using version %s, however version %s is available%s\n",
		colorStart, util.GetVersion(ctx), latestVersionSemver, colorEnd)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
askLoop:
	for {
		fmt.Println("Do you want to update [Y/n]: ")
		for !scanner.Scan() {
			break
		}
		if scanner.Err() != nil {
			util.GetLog(ctx).WithError(err).Error("scan failed")
			continue
		}
		switch strings.ToLower(scanner.Text()) {
		case "", "y":
			break askLoop
		case "n":
			return nil
		default:
			continue
		}
	}

	archive, err := downloader.LatestDownload()
	if err != nil {
		return err
	}
	defer archive.Close()

	pkg, err := unpack(archive)
	if err != nil {
		return err
	}
	defer pkg.Close()

	err = verifiedUpdate(pkg)
	if err != nil {
		return err
	}

	if restartAfter {
		gracefulRestart(ctx)
	}

	return nil
}
