package chlib

import (
	"bufio"
	"bytes"
	"chkit-v2/chlib/dbconfig"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

type TcpApiHandler struct {
	cfg      dbconfig.TcpApiConfig
	authForm map[string]string
	socket   net.Conn
}

type TcpApiResult map[string]interface{}

func NewTcpApiHandler(cfg dbconfig.TcpApiConfig, uuid, token string) *TcpApiHandler {
	handler := &TcpApiHandler{cfg: cfg}
	handler.authForm = make(map[string]string)
	handler.authForm["channel"] = uuid
	handler.authForm["token"] = token
	return handler
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
		return nil
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
