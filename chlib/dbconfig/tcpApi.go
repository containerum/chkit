package dbconfig

import (
	"chkit-v2/helpers"
	"fmt"
)

type TcpApiConfig struct {
	Address    string `mapconv:"address"`
	BufferSize int    `mapconv:"buffersize"`
}

const tcpApiBucket = "tcpApi"

func init() {
	cfg := TcpApiConfig{
		Address:    DefaultTCPServer,
		BufferSize: DefaultBufferSize,
	}
	initializers[tcpApiBucket] = helpers.StructToMap(cfg)
}

func (d *ConfigDB) GetTcpApiConfig() (cfg TcpApiConfig, err error) {
	m, err := d.readTransactional(tcpApiBucket)
	if err != nil {
		return cfg, fmt.Errorf("load tcp api config: %s", err)
	}
	err = helpers.FillStruct(&cfg, m)
	if err != nil {
		return cfg, fmt.Errorf("fill tcp api config: %s", err)
	}
	return cfg, nil
}

func (d *ConfigDB) UpdateTcpApiConfig(cfg TcpApiConfig) error {
	return d.writeTransactional(helpers.StructToMap(cfg), tcpApiBucket)
}
