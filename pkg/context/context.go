package context

import (
	"path"

	"github.com/blang/semver"
	"github.com/containerum/chkit/pkg/client"
	"github.com/containerum/chkit/pkg/configdir"
	"github.com/containerum/chkit/pkg/model"
)

var Context = struct {
	Version      string
	ConfigPath   string
	ConfigDir    string
	APIaddr      string
	Fingerprint  string
	Namespace    string
	Tokens       model.Tokens
	ClientConfig model.Config
	Client       *chClient.Client
}{
	Version:    semver.MustParse("3.0.1-alpha").String(),
	ConfigDir:  configdir.ConfigDir(),
	ConfigPath: path.Join(configdir.ConfigDir(), "config.toml"),
}
