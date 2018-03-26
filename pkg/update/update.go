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

func AskForUpdate(ctx *cli.Context, latestVersion semver.Version) (bool, error) {
	// check if we have terminal supports escape sequences
	var colorStart, colorEnd string
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		colorStart = "\x1b[31;1m"
		colorEnd = "\x1b[0m"
	}
	fmt.Printf("%sYou are using version %s, however version %s is available%s\n",
		colorStart, util.GetVersion(ctx), latestVersion, colorEnd)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	for {
		fmt.Println("Do you want to update [Y/n]: ")
		for !scanner.Scan() {
			break
		}
		if err := scanner.Err(); err != nil {
			util.GetLog(ctx).WithError(err).Error("scan failed")
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

func Update(ctx *cli.Context, downloader LatestCheckerDownloader, restartAfter bool) error {
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
