package dbconfig

import (
	"fmt"
	"time"

	"github.com/containerum/chkit/helpers"
)

type HttpApiConfig struct {
	Server  string        `mapconv:"server"`
	Timeout time.Duration `mapconv:"timeout"`
}

const httpApiBucket = "httpApi"

func init() {
	cfg := HttpApiConfig{
		Server:  DefaultHTTPServer,
		Timeout: DefaultHTTPTimeout,
	}
	initializers[httpApiBucket] = helpers.StructToMap(cfg)
}

func (d *ConfigDB) GetHttpApiConfig() (cfg HttpApiConfig, err error) {
	m, err := d.readTransactional(httpApiBucket)
	if err != nil {
		return cfg, fmt.Errorf("http api bucket get: %s", err)
	}
	err = helpers.FillStruct(&cfg, m)
	if err != nil {
		return cfg, fmt.Errorf("http api fill: %s", err)
	}
	return
}

func (d *ConfigDB) UpdateHttpApiConfig(cfg HttpApiConfig) error {
	return d.writeTransactional(helpers.StructToMap(cfg), httpApiBucket)
}
