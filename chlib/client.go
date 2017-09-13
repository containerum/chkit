package chlib

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/containerum/chkit.v2/chlib/dbconfig"

	jww "github.com/spf13/jwalterweatherman"
)

type Client struct {
	path          string
	version       string
	apiHandler    *HttpApiHandler
	tcpApiHandler *TcpApiHandler
	userConfig    *dbconfig.UserInfo
}

type GenericJson map[string]interface{}

func NewClient(db *dbconfig.ConfigDB, version, uuid string, np *jww.Notepad) (*Client, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	cfg, err := db.GetHttpApiConfig()
	if err != nil {
		return nil, err
	}
	tcpApiCfg, err := db.GetTcpApiConfig()
	if err != nil {
		return nil, err
	}
	client := &Client{
		path:    cwd,
		version: version,
	}
	userCfg, err := db.GetUserInfo()
	if err != nil {
		return nil, err
	}
	client.apiHandler = NewHttpApiHandler(cfg, uuid, userCfg.Token, np)
	client.tcpApiHandler = NewTcpApiHandler(tcpApiCfg, uuid, userCfg.Token, np)
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

func (c *Client) filterKind(apiResult TcpApiResult, kind string) (newApiResult TcpApiResult) {
	newApiResult = make(TcpApiResult)
	newResult := []interface{}{}
	for _, r := range apiResult["results"].([]interface{}) {
		if r.(map[string]interface{})["data"].(map[string]interface{})["kind"].(string) == kind {
			newResult = append(newResult, r)
		}
	}
	newApiResult["results"] = newResult
	for k, v := range apiResult {
		if k != "results" {
			newApiResult[k] = v
		}
	}
	return
}

func (c *Client) Get(kind, name, nameSpace string) (apiResult TcpApiResult, err error) {
	if c.userConfig.Token == "" {
		return nil, fmt.Errorf("Token is empty. Please login or set it manually (see help for \"config\" command)")
	}
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
	if err != nil {
		return
	}
	tmpApiResult, err := c.tcpApiHandler.Receive()
	// remove invalid kinds
	if kind == KindNamespaces && name == "" {
		apiResult = c.filterKind(tmpApiResult, "ResourceQuota")
	} else {
		apiResult = tmpApiResult
	}
	return
}

func (c *Client) Set(deploy, container, parameter, value, nameSpace string) (res TcpApiResult, err error) {
	if c.userConfig.Token == "" {
		return nil, fmt.Errorf("Token is empty. Please login or set it manually (see help for \"config\" command)")
	}
	_, err = c.tcpApiHandler.Connect()
	if err != nil {
		return
	}
	defer c.tcpApiHandler.Close()
	if nameSpace == "" {
		nameSpace = c.userConfig.Namespace
	}
	var httpResult HttpApiResult
	if container != "" {
		req := GenericJson{
			"name":    deploy,
			parameter: value,
		}
		httpResult, err = c.apiHandler.SetForContainer(req, container, nameSpace)
	} else {
		req := make(GenericJson)
		switch parameter {
		case "replicas":
			replicas, err := strconv.Atoi(value)
			if err != nil || replicas <= 0 {
				return res, fmt.Errorf("invalid replicas count")
			}
			req[parameter] = replicas
		default:
			req[parameter] = value
		}
		httpResult, err = c.apiHandler.SetForDeploy(req, deploy, nameSpace)
	}
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

func (c *Client) Create(jsonToSend GenericJson) (res TcpApiResult, err error) {
	if c.userConfig.Token == "" {
		return nil, fmt.Errorf("Token is empty. Please login or set it manually (see help for \"config\" command)")
	}
	_, err = c.tcpApiHandler.Connect()
	if err != nil {
		return
	}
	defer c.tcpApiHandler.Close()
	metaDataI, hasMd := jsonToSend["metadata"]
	if !hasMd {
		return res, fmt.Errorf("JSON must have \"metadata\" parameter")
	}
	metaData, validMd := metaDataI.(map[string]interface{})
	if !validMd {
		return res, fmt.Errorf("metadata must be object")
	}
	nameSpaceI, hasNs := metaData["namespace"]
	var nameSpace string
	if hasNs {
		var valid bool
		nameSpace, valid = nameSpaceI.(string)
		if !valid {
			return res, fmt.Errorf("namespace must be string")
		}
	} else {
		nameSpace = c.userConfig.Namespace
	}
	kindI, hasKind := jsonToSend["kind"]
	if !hasKind {
		return res, fmt.Errorf("JSON must have kind field")
	}
	kind, valid := kindI.(string)
	if !valid {
		return res, fmt.Errorf("kind must be string")
	}
	httpResult, err := c.apiHandler.Create(jsonToSend, kind, nameSpace)
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

func (c *Client) Delete(kind, name, nameSpace string, allPods bool) (res TcpApiResult, err error) {
	if c.userConfig.Token == "" {
		return nil, fmt.Errorf("Token is empty. Please login or set it manually (see help for \"config\" command)")
	}
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
		httpResult, err = c.apiHandler.Delete(kind, name, nameSpace, allPods)
	} else {
		httpResult, err = c.apiHandler.DeleteNameSpaces(name)
	}
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

func (c *Client) constructExpose(name string, ports []Port, nameSpace string) (ret GenericJson, err error) {
	req := new(Service)
	req.Kind = "Service"

	nsHash := sha256.Sum256([]byte(nameSpace))
	nameLabelKey := hex.EncodeToString(nsHash[:])[:32]
	req.Spec.Selector = map[string]string{
		nameLabelKey: name,
	}

	external := "true"
	for _, port := range ports {
		if port.Port != 0 {
			external = "false"
		}
	}
	req.Spec.Ports = ports
	req.Metadata.Labels = map[string]string{
		nameLabelKey: name,
		"external":   external,
	}

	nameHash := md5.Sum([]byte(name + time.Now().Format("2006-01-02 15:04:05.000000")))
	_, err = c.Get(KindDeployments, name, nameSpace)
	if err != nil {
		return nil, fmt.Errorf("expose construct: %s", err)
	}
	req.Metadata.Name = fmt.Sprintf("%s-%s", name, hex.EncodeToString(nameHash[:])[:4])

	b, _ := json.MarshalIndent(req, "", "    ")
	err = ioutil.WriteFile(ExposeFile, b, 0600)
	if err != nil {
		return nil, fmt.Errorf("expose write file: %s", err)
	}
	err = json.Unmarshal(b, &ret)
	return
}

func (c *Client) Expose(name string, ports []Port, nameSpace string) (res TcpApiResult, err error) {
	if c.userConfig.Token == "" {
		return nil, fmt.Errorf("Token is empty. Please login or set it manually (see help for \"config\" command)")
	}
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
	httpResult, err := c.apiHandler.Expose(req, nameSpace)
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

type ConfigureParams struct {
	Image    string
	Ports    []int
	Labels   map[string]string
	Env      []EnvVar
	CPU      string
	Memory   string
	Replicas int
	Command  []string
}

func (c *Client) constructRun(name string, params ConfigureParams) (ret GenericJson, err error) {
	req := new(Deploy)
	req.Kind = "Deployment"
	req.Metadata.Name = name
	req.Metadata.Labels = params.Labels
	req.Spec.Replicas = params.Replicas
	req.Spec.Template.Metadata.Name = name
	req.Spec.Template.Metadata.Labels = params.Labels
	containers := make([]Container, 1)
	containers[0].Name = name
	containers[0].Image = params.Image
	if len(params.Ports) != 0 {
		for _, p := range params.Ports {
			containers[0].Ports = append(containers[0].Ports, Port{ContainerPort: p})
		}
	}
	containers[0].Command = params.Command
	containers[0].Env = params.Env
	containers[0].Resources.Requests = &HwResources{CPU: params.CPU, Memory: params.Memory}
	req.Spec.Template.Spec.Containers = containers
	b, _ := json.MarshalIndent(req, "", "    ")
	err = ioutil.WriteFile(RunFile, b, 0600)
	if err != nil {
		return nil, fmt.Errorf("run write file: %s", err)
	}
	err = json.Unmarshal(b, &ret)
	return
}

func (c *Client) Run(name string, params ConfigureParams, nameSpace string) (res TcpApiResult, err error) {
	if c.userConfig.Token == "" {
		return nil, fmt.Errorf("Token is empty. Please login or set it manually (see help for \"config\" command)")
	}
	_, err = c.tcpApiHandler.Connect()
	if err != nil {
		return
	}
	defer c.tcpApiHandler.Close()
	req, err := c.constructRun(name, params)
	if err != nil {
		return
	}
	if nameSpace == "" {
		nameSpace = c.userConfig.Namespace
	}
	httpResult, err := c.apiHandler.Run(req, nameSpace)
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

func (c *Client) GetVolume(name string) (res interface{}, err error) {
	if c.userConfig.Token == "" {
		return nil, fmt.Errorf("Token is empty. Please login or set it manually (see help for \"config\" command)")
	}
	res, err = c.apiHandler.GetVolume(name)
	if tmp, ok := res.(map[string]interface{}); ok {
		if _, ok := tmp["label"]; !ok {
			return res, fmt.Errorf("volume %s was not found", name)
		}
	}
	return
}
