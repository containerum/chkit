package chlib

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type Client struct {
	path          string
	version       string
	apiHandler    *HttpApiHandler
	tcpApiHandler *TcpApiHandler
	userConfig    *UserInfo
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
	userCfg, err := GetUserInfo()
	if err != nil {
		return nil, err
	}
	tcpApiCfg.Token = userCfg.Token
	cfg.Token = userCfg.Token
	client.apiHandler = NewHttpApiHandler(cfg)
	client.tcpApiHandler = NewTcpApiHandler(tcpApiCfg)
	client.userConfig = &userCfg
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
	err = apiResult.HandleApiResult()
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

func (c *Client) Get(kind, name, nameSpace string) (apiResult TcpApiResult, err error) {
	_, err = c.tcpApiHandler.Connect()
	if err != nil {
		return
	}
	defer c.tcpApiHandler.Close()
	var httpResult HttpApiResult
	if nameSpace == "" {
		nameSpace = c.userConfig.Namespace
	}
	if kind != KindNamespaces {
		httpResult, err = c.apiHandler.Get(kind, name, nameSpace)
	} else {
		httpResult, err = c.apiHandler.GetNameSpaces(name)
	}
	if err != nil {
		return
	}
	err = httpResult.HandleApiResult()
	apiResult, err = c.tcpApiHandler.Receive()
	return
}

func (c *Client) Set(field, container, value, nameSpace string) (res TcpApiResult, err error) {
	_, err = c.tcpApiHandler.Connect()
	if err != nil {
		return
	}
	if nameSpace == "" {
		nameSpace = c.userConfig.Namespace
	}
	defer c.tcpApiHandler.Close()
	reqData := GenericJson{"name": container}
	reqData[field] = value
	httpResult, err := c.apiHandler.Set(reqData, container, nameSpace)
	if err != nil {
		return
	}
	err = httpResult.HandleApiResult()
	if err != nil {
		return
	}
	res, err = c.tcpApiHandler.Receive()
	if err != nil {
		return
	}
	err = res.CheckHttpStatus()
	return
}

func (c *Client) Create(jsonToSend GenericJson) (apiResult TcpApiResult, err error) {
	metaDataI, hasMd := jsonToSend["metadata"]
	if !hasMd {
		return apiResult, fmt.Errorf("JSON must have \"metadata\" parameter")
	}
	metaData, validMd := metaDataI.(map[string]interface{})
	if !validMd {
		return apiResult, fmt.Errorf("metadata must be object")
	}
	nameSpaceI, hasNs := metaData["namespace"]
	var nameSpace string
	if hasNs {
		var valid bool
		nameSpace, valid = nameSpaceI.(string)
		if !valid {
			return apiResult, fmt.Errorf("namespace must be string")
		}
	} else {
		nameSpace = c.userConfig.Namespace
	}
	kindI, hasKind := jsonToSend["kind"]
	if !hasKind {
		return apiResult, fmt.Errorf("JSON must have kind field")
	}
	kind, valid := kindI.(string)
	if !valid {
		return apiResult, fmt.Errorf("kind must be string")
	}
	httpResult, err := c.apiHandler.Create(jsonToSend, kind, nameSpace)
	if err != nil {
		return
	}
	err = httpResult.HandleApiResult()
	return
}

func (c *Client) Delete(kind, name, nameSpace string, allPods bool) (err error) {
	var httpResult HttpApiResult
	if kind != KindNamespaces {
		httpResult, err = c.apiHandler.Delete(kind, name, nameSpace, allPods)
	} else {
		httpResult, err = c.apiHandler.DeleteNameSpaces(name)
	}
	if err != nil {
		return
	}
	err = httpResult.HandleApiResult()
	return
}

func (c *Client) constructExpose(name string, ports []Port, nameSpace string) (ret GenericJson, err error) {
	labels := make(map[string]string)
	labels["external"] = "true"
	nsHash := sha256.Sum256([]byte(nameSpace))
	labels[hex.EncodeToString(nsHash[:])[:32]] = nameSpace
	nameHash := md5.Sum([]byte(name + time.Now().Format("2006-01-02 15:04:05.000000")))
	_, err = c.Get(KindDeployments, name, nameSpace)
	if err != nil {
		return nil, fmt.Errorf("expose construct: %s", err)
	}
	req := new(Service)
	req.Spec.Ports = ports
	req.Metadata.Labels = labels
	req.Metadata.Name = fmt.Sprintf("%s-%s", name, hex.EncodeToString(nameHash[:])[:4])
	req.Spec.Selector = labels
	exposeJsonPath := path.Join(configPath, srcFolder, templatesFolder, exposeFile)
	b, _ := json.MarshalIndent(req, "", "    ")
	err = ioutil.WriteFile(exposeJsonPath, b, 0600)
	if err != nil {
		return nil, fmt.Errorf("expose write file: %s", err)
	}
	err = json.Unmarshal(b, &ret)
	ret["kind"] = "Service"
	return
}

func (c *Client) Expose(name string, ports []Port, nameSpace string) (apiResult TcpApiResult, err error) {
	if nameSpace == "" {
		nameSpace = c.userConfig.Namespace
	}
	var req GenericJson
	req, err = c.constructExpose(name, ports, nameSpace)
	if err != nil {
		return
	}
	_, err = c.tcpApiHandler.Connect()
	if err != nil {
		return
	}
	defer c.tcpApiHandler.Close()
	var httpResult HttpApiResult
	httpResult, err = c.apiHandler.Expose(req, nameSpace)
	if err != nil {
		return
	}
	err = httpResult.HandleApiResult()
	if err != nil {
		return
	}
	apiResult, err = c.tcpApiHandler.Receive()
	return
}
