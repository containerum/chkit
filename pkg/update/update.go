package update

import (
	"crypto"
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"time"
	"unicode/utf8"

	"math"

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

func speedDecorator(minWidth int, conf byte) decor.DecoratorFunc {
	format := "%%"
	if (conf & decor.DidentRight) != 0 {
		format += "-"
	}
	format += "%ds"
	var (
		prevCount         int64
		prevDuration      time.Duration
		integrSpeed       float64
		measurementsCount float64
	)
	return func(s *decor.Statistics, widthAccumulator chan<- int, widthDistributor <-chan int) string {
		var str string
		measurementsCount++
		curSpeed := float64(s.Current-prevCount) / (s.TimeElapsed - prevDuration).Seconds() // bytes per second
		if !math.IsNaN(curSpeed) && !math.IsInf(curSpeed, 1) && !math.IsInf(curSpeed, -1) {
			integrSpeed += curSpeed
		}
		speedToShow := integrSpeed / measurementsCount
		prevCount = s.Current
		prevDuration = s.TimeElapsed
		switch {
		case speedToShow < 1<<10:
			str = fmt.Sprintf("%.1f B/s", speedToShow)
		case speedToShow >= 1<<10 && speedToShow < 1<<20:
			str = fmt.Sprintf("%.1f KiB/s", speedToShow/(1<<10))
		case speedToShow >= 1<<20 && speedToShow < 1<<30:
			str = fmt.Sprintf("%.1f MiB/s", speedToShow/(1<<20))
		default:
			str = fmt.Sprintf("%.1f GiB/s", speedToShow/(1<<30))
		}
		if (conf & decor.DwidthSync) != 0 {
			widthAccumulator <- utf8.RuneCountInString(str)
			max := <-widthDistributor
			if (conf & decor.DextraSpace) != 0 {
				max++
			}
			return fmt.Sprintf(fmt.Sprintf(format, max), str)
		}
		return fmt.Sprintf(fmt.Sprintf(format, minWidth), str)
	}
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
		decor.CountersKiloByte("%.1f / %.1f", 3, 0),
		decor.StaticName("(", 1, 0),
		decor.Percentage(3, 0),
		decor.StaticName(")", 1, 0),
	), mpb.AppendDecorators(
		speedDecorator(3, decor.DextraSpace),
		decor.StaticName(" ETA:", 3, decor.DextraSpace),
		decor.ETA(3, decor.DextraSpace),
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
