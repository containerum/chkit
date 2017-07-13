package chlib

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/kfeofantov/chkit-v2/helpers"
)

type TcpApiConfig struct {
	Address    net.IP `mapconv:"address"`
	Port       int    `mapconv:"port"`
	BufferSize int    `mapconv:"buffersize"`
	Uuid       string `mapconv:"-"`
	Token      string `mapconv:"-"`
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

type TcpApiHandler struct {
	cfg      TcpApiConfig
	authForm map[string]string
	socket   net.Conn
}

type TcpApiResult map[string]interface{}

func GetTcpApiConfig() (cfg TcpApiConfig, err error) {
	m, err := readFromBucket(tcpApiBucket)
	if err != nil {
		return cfg, fmt.Errorf("load tcp api config: %s", err)
	}
	err = helpers.FillStruct(&cfg, m)
	if err != nil {
		return cfg, fmt.Errorf("fill tcp api config: %s", err)
	}
	return cfg, nil
}

func UpdateTcpApiConfig(cfg TcpApiConfig) error {
	return pushToBucket(tcpApiBucket, helpers.StructToMap(cfg))
}

func NewTcpApiHandler(cfg TcpApiConfig) *TcpApiHandler {
	return &TcpApiHandler{
		cfg: cfg,
		authForm: map[string]string{
			"channel": cfg.Uuid,
			"token":   cfg.Token,
		},
	}
}

func (t *TcpApiHandler) Connect() (result TcpApiResult, err error) {
	t.socket, err = net.Dial("tcp", net.JoinHostPort(t.cfg.Address.String(), strconv.Itoa(t.cfg.Port)))
	if err != nil {
		return result, fmt.Errorf("tcp connect: %s", err)
	}
	hello, err := json.Marshal(t.authForm)
	if err != nil {
		return result, fmt.Errorf("authForm encode: %s", err)
	}
	hello = append(hello, '\n')
	_, err = t.socket.Write(hello)
	if err != nil {
		return result, fmt.Errorf("hello send: %s", err)
	}
	str, err := bufio.NewReader(t.socket).ReadSlice('\n')
	if err != nil {
		return result, fmt.Errorf("hello receive: %s", err)
	}
	err = json.Unmarshal(str, &result)
	if err != nil {
		return result, fmt.Errorf("hello decode: %s", err)
	}
	return result, result.CheckHttpStatus()
}

func (t *TcpApiHandler) Receive() (result TcpApiResult, err error) {
	var data []byte
	for buf := make([]byte, t.cfg.BufferSize); !bytes.ContainsRune(buf, '\n'); {
		n, err := t.socket.Read(buf)
		if err != nil {
			return result, err
		}
		data = append(data, buf[:n]...)
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return result, err
	}
	return result, result.CheckHttpStatus()
}

func (t *TcpApiHandler) Close() {
	if t.socket != nil {
		t.socket.Close()
	}
}

func (t *TcpApiResult) CheckHttpStatus() error {
	if t == nil {
		return fmt.Errorf("TCP API result is nil")
	}
	errContent, hasError := (*t)["error"]
	if !hasError {
		return nil
	}
	errStr, isStr := errContent.(string)
	if !isStr {
		fmt.Errorf("Got non-string error")
	}
	return fmt.Errorf(errStr)
}
