package client

import (
	"encoding/json"
	"io"

	"github.com/containerum/chkit/pkg/err"
)

var (
	ErrUnableToPackTokens    = err.Err("unable to pack tokens")
	ErrUnableToSaveTokenFile = err.Err("unable to save token file")
)

func (client *ChkitClient) SaveTokens(wr io.Writer) error {
	packedTokens, err := json.Marshal(client.tokens)
	if err != nil {
		return ErrUnableToPackTokens.Wrap(err)
	}
	_, err = wr.Write(packedTokens)
	if err != nil {
		return ErrUnableToSaveTokenFile.Wrap(err)
	}
	return nil
}
