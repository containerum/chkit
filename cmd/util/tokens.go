package util

import (
	"encoding/json"
	"os"
	"path"

	. "github.com/containerum/chkit/cmd/context"
	"github.com/containerum/chkit/pkg/model"
)

// SaveTokens -- save tokens in config path
func SaveTokens(tokens model.Tokens) error {
	file, err := os.Create(path.Join(Context.ConfigDir, "tokens"))
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(tokens)
}

// LoadTokens -- loads tokens from fs
func LoadTokens() (model.Tokens, error) {
	tokens := model.Tokens{}
	file, err := os.Open(path.Join(Context.ConfigPath, "tokens"))
	if err != nil && !os.IsNotExist(err) {
		return tokens, err
	} else if err != nil && os.IsNotExist(err) {
		return tokens, nil
	}
	return tokens, json.NewDecoder(file).Decode(&tokens)
}
