package chlib

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"

	"github.com/containerum/chkit/chlib/dbconfig"

	"time"

	jww "github.com/spf13/jwalterweatherman"
)

type TcpApiHandler struct {
	Config   *dbconfig.TcpApiConfig
	UserInfo *dbconfig.UserInfo
	Channel  string
	Np       *jww.Notepad

	socket net.Conn
}

type TcpApiResult map[string]interface{}

func (t *TcpApiHandler) Connect() (result TcpApiResult, err error) {
	t.Np.SetPrefix("TCP")
	t.Np.DEBUG.Println("connect", t.Config.Address)
	t.socket, err = net.Dial("tcp", t.Config.Address)
	if err != nil {
		return result, fmt.Errorf("tcp connect: %s", err)
	}
	if err := t.socket.SetReadDeadline(time.Now().Add(t.Config.Timeout)); err != nil {
		return result, fmt.Errorf("tcp set read deadline: %s", err)
	}
	if err := t.socket.SetWriteDeadline(time.Now().Add(t.Config.Timeout)); err != nil {
		return result, fmt.Errorf("tcp set write deadline: %s", err)
	}
	hello, err := json.Marshal(map[string]string{
		"channel": t.Channel,
		"token":   t.UserInfo.Token,
	})
	if err != nil {
		return result, fmt.Errorf("authForm encode: %s", err)
	}
	hello = append(hello, '\n')
	t.Np.DEBUG.Println("TCP Auth")
	n, err := t.socket.Write(hello)
	t.Np.DEBUG.Printf("Write %d bytes\n", n)
	if err != nil {
		return result, fmt.Errorf("hello send: %s", err)
	}
	t.Np.DEBUG.Println("Auth response read")
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
	t.Np.SetPrefix("TCP")
	var data []byte
	for buf := make([]byte, t.Config.BufferSize); !bytes.ContainsRune(buf, '\n'); {
		n, err := t.socket.Read(buf)
		if err != nil {
			return result, err
		}
		data = append(data, buf[:n]...)
	}
	t.Np.DEBUG.Printf("Received %d bytes\n", len(data))
	err = json.Unmarshal(data, &result)
	if err != nil {
		return result, err
	}
	cmdId, hasCmdId := result["id"]
	if hasCmdId {
		t.Np.DEBUG.Printf("Command ID: %v\n", cmdId)
	}
	return result, result.CheckHttpStatus()
}

func (t *TcpApiHandler) Close() {
	if t.socket != nil {
		t.Np.DEBUG.Println("Close connection")
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
		fmt.Errorf("got non-string error")
	}
	return fmt.Errorf(errStr)
}
