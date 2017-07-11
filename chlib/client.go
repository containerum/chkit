package chlib

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
)

type Client struct {
	path          string
	version       string
	apiHandler    *HttpApiHandler
	tcpApiHandler *TcpApiHandler
}

type GenericJson map[string]interface{}

func NewClient(version, uuid string) (*Client, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	cfg, err := GetHttpApiCfg()
	if err != nil {
		return nil, err
	}
	cfg.Uuid = uuid
	tcpApiCfg, err := GetTcpApiConfig()
	if err != nil {
		return nil, err
	}
	tcpApiCfg.Uuid = uuid
	client := &Client{
		path:    cwd,
		version: version,
	}
	client.apiHandler = NewHttpApiHandler(cfg)
	client.tcpApiHandler = NewTcpApiHandler(tcpApiCfg)
	return client, nil
}

func (c *Client) Login(login, password string) (token string, err error) {
	passwordHash := md5.Sum([]byte(login + password))
	jsonToSend := GenericJson{
		"username": login,
		"password": hex.EncodeToString(passwordHash[:]),
	}
	apiResult, err := c.apiHandler.Login(jsonToSend)
	if err != nil {
		return "", err
	}
	err = c.handleApiResult(apiResult)
	if err != nil {
		return "", err
	}
	tokenI, hasToken := apiResult["token"]
	if !hasToken {
		return "", fmt.Errorf("api result don`t have token")
	}
	token, isString := tokenI.(string)
	if !isString {
		return "", fmt.Errorf("received non-string token")
	}
	return token, nil
}

func (c *Client) handleApiResult(apiResult HttpApiResult) error {
	errCont, hasErr := apiResult["error"]
	if hasErr {
		return fmt.Errorf("api error: %v", errCont)
	}
	return nil
}
