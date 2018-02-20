package client

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/containerum/chkit/pkg/err"
)

var (
	ErrUnableToPackTokens    = err.Err("unable to pack tokens")
	ErrUnableToSaveTokenFile = err.Err("unable to save token file")
)

func (client *ChkitClient) SaveTokens() error {
	packedTokens, err := json.Marshal(client.tokens)
	if err != nil {
		return ErrUnableToPackTokens.Wrap(err)
	}
	err = ioutil.WriteFile(client.config.TokenFile,
		packedTokens,
		os.ModePerm)
	if err != nil {
		return ErrUnableToSaveTokenFile.Wrap(err)
	}
	return nil
}
