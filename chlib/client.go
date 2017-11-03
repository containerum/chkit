package chlib

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"strings"

	"github.com/containerum/chkit/chlib/dbconfig"
	"github.com/containerum/solutions"
)

type Client struct {
	ApiHandler    *HttpApiHandler
	TcpApiHandler *TcpApiHandler
	UserConfig    *dbconfig.UserInfo
}

type GenericJson map[string]interface{}

func (c *Client) Login(login, password string) (err error) {
	passwordHash := md5.Sum([]byte(login + password))
	jsonToSend := GenericJson{
		"username": login,
		"password": hex.EncodeToString(passwordHash[:]),
	}
	apiResult, err := c.ApiHandler.Login(jsonToSend)
	if err != nil {
		return
	}
	err = apiResult.HandleApiResult()
	if err != nil {
		return
	}
	tokenI, hasToken := apiResult["token"]
	if !hasToken {
		return fmt.Errorf("api result don`t have token")
	}
	token, isString := tokenI.(string)
	if !isString {
		return fmt.Errorf("received non-string token")
	}
	c.UserConfig.Token = token
	return nil
}

func (c *Client) Get(kind, name, nameSpace string) (apiResult TcpApiResult, err error) {
	if c.UserConfig.Token == "" {
		return nil, fmt.Errorf("Token is empty. Please login or set it manually (see help for \"config\" command)")
	}
	_, err = c.TcpApiHandler.Connect()
	if err != nil {
		return
	}
	defer c.TcpApiHandler.Close()
	var httpResult HttpApiResult
	if nameSpace == "" {
		nameSpace = c.UserConfig.Namespace
	}
	if kind != KindNamespaces {
		httpResult, err = c.ApiHandler.Get(kind, name, nameSpace)
	} else {
		httpResult, err = c.ApiHandler.GetNameSpaces(name)
	}
	if err != nil {
		return
	}
	err = httpResult.HandleApiResult()
	if err != nil {
		return
	}
	return c.TcpApiHandler.Receive()
}

func (c *Client) Set(deploy, container, parameter, value, nameSpace string) (res TcpApiResult, err error) {
	if c.UserConfig.Token == "" {
		return nil, fmt.Errorf("Token is empty. Please login or set it manually (see help for \"config\" command)")
	}
	_, err = c.TcpApiHandler.Connect()
	if err != nil {
		return
	}
	defer c.TcpApiHandler.Close()
	if nameSpace == "" {
		nameSpace = c.UserConfig.Namespace
	}
	var httpResult HttpApiResult
	if container != "" {
		req := GenericJson{
			"name":    deploy,
			parameter: value,
		}
		httpResult, err = c.ApiHandler.SetForContainer(req, container, nameSpace)
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
		httpResult, err = c.ApiHandler.SetForDeploy(req, deploy, nameSpace)
	}
	if err != nil {
		return
	}
	err = httpResult.HandleApiResult()
	if err != nil {
		return
	}
	res, err = c.TcpApiHandler.Receive()
	if err != nil {
		return
	}
	err = res.CheckHttpStatus()
	return
}

func (c *Client) Create(jsonToSend GenericJson) (res TcpApiResult, err error) {
	if c.UserConfig.Token == "" {
		return nil, fmt.Errorf("Token is empty. Please login or set it manually (see help for \"config\" command)")
	}
	_, err = c.TcpApiHandler.Connect()
	if err != nil {
		return
	}
	defer c.TcpApiHandler.Close()
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
		nameSpace = c.UserConfig.Namespace
	}
	kindI, hasKind := jsonToSend["kind"]
	if !hasKind {
		return res, fmt.Errorf("JSON must have kind field")
	}
	kind, valid := kindI.(string)
	if !valid {
		return res, fmt.Errorf("kind must be string")
	}
	httpResult, err := c.ApiHandler.Create(jsonToSend, kind, nameSpace)
	if err != nil {
		return
	}
	err = httpResult.HandleApiResult()
	if err != nil {
		return
	}
	res, err = c.TcpApiHandler.Receive()
	if err != nil {
		return
	}
	err = res.CheckHttpStatus()
	return
}

func (c *Client) Delete(kind, name, nameSpace string, allPods bool) (res TcpApiResult, err error) {
	if c.UserConfig.Token == "" {
		return nil, fmt.Errorf("Token is empty. Please login or set it manually (see help for \"config\" command)")
	}
	_, err = c.TcpApiHandler.Connect()
	if err != nil {
		return
	}
	defer c.TcpApiHandler.Close()
	var httpResult HttpApiResult
	if nameSpace == "" {
		nameSpace = c.UserConfig.Namespace
	}
	if kind != KindNamespaces {
		httpResult, err = c.ApiHandler.Delete(kind, name, nameSpace, allPods)
	} else {
		httpResult, err = c.ApiHandler.DeleteNameSpaces(name)
	}
	if err != nil {
		return
	}
	err = httpResult.HandleApiResult()
	if err != nil {
		return
	}
	res, err = c.TcpApiHandler.Receive()
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
	if c.UserConfig.Token == "" {
		return nil, fmt.Errorf("Token is empty. Please login or set it manually (see help for \"config\" command)")
	}
	if nameSpace == "" {
		nameSpace = c.UserConfig.Namespace
	}
	var req GenericJson
	req, err = c.constructExpose(name, ports, nameSpace)
	if err != nil {
		return
	}
	_, err = c.TcpApiHandler.Connect()
	if err != nil {
		return
	}
	defer c.TcpApiHandler.Close()
	httpResult, err := c.ApiHandler.Expose(req, nameSpace)
	if err != nil {
		return
	}
	err = httpResult.HandleApiResult()
	if err != nil {
		return
	}
	res, err = c.TcpApiHandler.Receive()
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
	Volumes  []Volume
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
	containers[0].VolumeMounts = params.Volumes
	req.Spec.Template.Spec.Containers = containers
	volNames := make(map[string]bool)
	for _, v := range params.Volumes {
		volNames[v.Label] = true
	}
	for k := range volNames {
		req.Spec.Template.Spec.Volumes = append(req.Spec.Template.Spec.Volumes, VolumeName{
			Name: k,
		})
	}
	b, _ := json.MarshalIndent(req, "", "    ")
	err = ioutil.WriteFile(RunFile, b, 0600)
	if err != nil {
		return nil, fmt.Errorf("run write file: %s", err)
	}
	err = json.Unmarshal(b, &ret)
	return
}

func (c *Client) Run(name string, params ConfigureParams, nameSpace string) (res TcpApiResult, err error) {
	if c.UserConfig.Token == "" {
		return nil, fmt.Errorf("Token is empty. Please login or set it manually (see help for \"config\" command)")
	}
	_, err = c.TcpApiHandler.Connect()
	if err != nil {
		return
	}
	defer c.TcpApiHandler.Close()
	req, err := c.constructRun(name, params)
	if err != nil {
		return
	}
	if nameSpace == "" {
		nameSpace = c.UserConfig.Namespace
	}
	httpResult, err := c.ApiHandler.Run(req, nameSpace)
	if err != nil {
		return
	}
	err = httpResult.HandleApiResult()
	if err != nil {
		return
	}
	res, err = c.TcpApiHandler.Receive()
	if err != nil {
		return
	}
	err = res.CheckHttpStatus()
	return
}

func (c *Client) GetVolume(name string) (res interface{}, err error) {
	if c.UserConfig.Token == "" {
		return nil, fmt.Errorf("Token is empty. Please login or set it manually (see help for \"config\" command)")
	}
	res, err = c.ApiHandler.GetVolume(name)
	if err != nil {
		return
	}
	if tmp, ok := res.(map[string]interface{}); ok {
		if _, ok := tmp["label"]; !ok {
			return res, fmt.Errorf("volume %s was not found", name)
		}
	}
	return
}

func (c *Client) RunSolution(solutionDir string, envArgs []string, nameSpace string) error {
	if c.UserConfig.Token == "" {
		return fmt.Errorf("Token is empty. Please login or set it manually (see help for \"config\" command)")
	}

	solution, err := solutions.OpenSolution(solutionDir)
	if err != nil {
		return fmt.Errorf("open solution: %v", err)
	}

	if !envSliceValidate(envArgs) {
		return fmt.Errorf("invalid environment variable argument detected")
	}

	var env = make(map[string]interface{})
	for _, v := range envArgs {
		envVar := strings.Split(v, "=")
		env[envVar[0]] = envVar[1]
	}

	solution.AddValues(env)

	if nameSpace == "" {
		nameSpace = c.UserConfig.Namespace
	}

	seq, err := solution.GenerateRunSequence(nameSpace)
	if err != nil {
		return fmt.Errorf("generate solution run sequence: %v", err)
	}

	for _, v := range seq {
		var jsonToSend GenericJson
		if err := json.Unmarshal([]byte(v.Config), &jsonToSend); err != nil {
			return fmt.Errorf("config unmarshal: %v", err)
		}
		if _, err := c.ApiHandler.Create(jsonToSend, v.Type, nameSpace); err != nil {
			return fmt.Errorf("%s create: %v", v.Type, err)
		}
	}

	return nil
}
