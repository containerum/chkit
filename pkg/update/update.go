package update

import (
	"crypto"
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/blang/semver"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/inconshreveable/go-update"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
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

func Update(currentVersion semver.Version, downloader LatestCheckerDownloader, restartAfter bool) error {
	latestVersion, err := downloader.LatestVersion()
	if err != nil {
		return err
	}
	if latestVersion.LE(currentVersion) {
		fmt.Println("You already using latest version. Update not needed.")
		return nil
	}

	archive, size, err := downloader.Download(latestVersion)
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
	fmt.Printf("Updated to version %s\n", latestVersion)

	if restartAfter {
		gracefulRestart()
	}

	return nil
}
