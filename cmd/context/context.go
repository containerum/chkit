package context

import (
	"path"

	"github.com/blang/semver"
	"github.com/containerum/chkit/cmd/config_dir"
	"github.com/containerum/chkit/pkg/client"
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
	ConfigDir:  confDir.ConfigDir(),
	ConfigPath: path.Join(confDir.ConfigDir(), "config.toml"),
}
