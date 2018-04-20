package configuration

import (
	"encoding/json"
	"os"
	"path"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model"
)

// SaveTokens -- save tokens in config path
func SaveTokens(ctx *context.Context, tokens model.Tokens) error {
	file, err := os.Create(path.Join(ctx.ConfigDir, "tokens"))
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(tokens)
}

// LoadTokens -- loads tokens from fs
func LoadTokens(ctx *context.Context) (model.Tokens, error) {
	tokens := model.Tokens{}
	file, err := os.Open(path.Join(ctx.ConfigDir, "tokens"))
	if err != nil && !os.IsNotExist(err) {
		return tokens, err
	} else if err != nil && os.IsNotExist(err) {
		return tokens, nil
	}
	return tokens, json.NewDecoder(file).Decode(&tokens)
}
