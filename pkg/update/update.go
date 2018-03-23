package update

import (
	"io/ioutil"

	"crypto"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/inconshreveable/go-update"
)

var PublicKey = `
-----BEGIN PUBLIC KEY-----
MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEtrVmBxQvheRArXjg2vG1xIprWGuCyESx
MMY8pjmjepSy2kuz+nl9aFLqmr+rDNdYvEBqQaZrYMc6k29gjvoQnQ==
-----END PUBLIC KEY-----
`

const (
	ErrUpdateApply = chkitErrors.Err("update apply failed")
)

func verifiedUpdate(upd *Update) error {
	checksum, err := ioutil.ReadAll(upd.Hash)
	if err != nil {
		return chkitErrors.Wrap(ErrUpdateApply, err)
	}

	signature, err := ioutil.ReadAll(upd.Signature)
	if err != nil {
		return chkitErrors.Wrap(ErrUpdateApply, err)
	}

	err = update.Apply(upd.Binary, update.Options{
		Checksum:  checksum,
		Signature: signature,
		Verifier:  update.NewECDSAVerifier(),
		Hash:      crypto.SHA256,
		PublicKey: []byte(PublicKey),
	})
	if err != nil {
		return chkitErrors.Wrap(ErrUpdateApply, err)
	}

	return nil
}
