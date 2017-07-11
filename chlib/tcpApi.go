package chlib

import (
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
	var hello bytes.Buffer
	err = json.NewEncoder(&hello).Encode(t.authForm)
	if err != nil {
		return result, fmt.Errorf("authForm encode: %s", err)
	}
	hello.WriteRune('\n')
	_, err = hello.WriteTo(t.socket)
	if err != nil {
		return result, fmt.Errorf("hello send: %s", err)
	}
	recvBuf := bytes.NewBuffer(make([]byte, t.cfg.BufferSize))
	_, err = recvBuf.ReadFrom(t.socket)
	if err != nil {
		return result, fmt.Errorf("hello receive: %s", err)
	}
	err = json.NewDecoder(recvBuf).Decode(&result)
	if err != nil {
		return result, fmt.Errorf("hello decode: %s", err)
	}
	return result, nil
}

func (t *TcpApiHandler) Receive() (result TcpApiResult, err error) {
	var data bytes.Buffer
	for data.Bytes()[data.Cap()-1] != '\n' {
		received := bytes.NewBuffer(make([]byte, t.cfg.BufferSize))
		_, err := received.ReadFrom(t.socket)
		if err != nil {
			return result, fmt.Errorf("tcp receive: %s", err)
		}
		data.Write(received.Bytes())
	}
	err = json.NewDecoder(&data).Decode(&result)
	return result, err
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
