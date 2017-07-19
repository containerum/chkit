package dbconfig

import (
	"chkit-v2/helpers"
	"fmt"
	"net"
)

type TcpApiConfig struct {
	Address    net.IP `mapconv:"address"`
	Port       int    `mapconv:"port"`
	BufferSize int    `mapconv:"buffersize"`
}

const tcpApiBucket = "tcpApi"

func init() {
	cfg := TcpApiConfig{
		Address:    net.IPv4zero,
		Port:       0,
		BufferSize: 1024,
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
