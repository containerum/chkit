package chlib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/containerum/chkit/chlib/dbconfig"

	jww "github.com/spf13/jwalterweatherman"
)

type HttpApiHandler struct {
	Config   *dbconfig.HttpApiConfig
	UserInfo *dbconfig.UserInfo
	Np       *jww.Notepad
	Channel  string
}

type HttpApiResult map[string]interface{}

func (h *HttpApiHandler) makeRequestGeneric(url, method string, jsonToSend GenericJson, result interface{}) (err error) {
	h.Np.SetPrefix("HTTP")
	client := http.Client{Timeout: h.Config.Timeout}
	marshalled, _ := json.Marshal(jsonToSend)
	h.Np.DEBUG.Printf("%s %s\n", method, url)
	request, err := http.NewRequest(method, url, bytes.NewBuffer(marshalled))
	request.Header.Add("Channel", h.Channel)
	request.Header.Add("Authorization", h.UserInfo.Token)
	if err != nil {
		return fmt.Errorf("http request create error: %s", err)
	}
	resp, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("http request execute error: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got non-ok http response: %s", resp.Status)
	}
	h.Np.DEBUG.Println("Result", resp.Status)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, result)
	if err != nil {
		err = fmt.Errorf("received error: %s", body)
	}
	return err
}

func (h *HttpApiHandler) makeRequest(url, method string, jsonToSend GenericJson) (result HttpApiResult, err error) {
	err = h.makeRequestGeneric(url, method, jsonToSend, &result)
	return
}

func (h *HttpApiHandler) Create(jsonToSend GenericJson, kind, nameSpace string) (result HttpApiResult, err error) {
	kind = fmt.Sprintf("%ss", strings.ToLower(kind))
	url := fmt.Sprintf("%s/namespaces/%s/%s", h.Config.Server, nameSpace, kind)
	return h.makeRequest(url, http.MethodPost, jsonToSend)
}

func (h *HttpApiHandler) Expose(jsonToSend GenericJson, nameSpace string) (result HttpApiResult, err error) {
	url := fmt.Sprintf("%s/namespaces/%s/services", h.Config.Server, nameSpace)
	return h.makeRequest(url, http.MethodPost, jsonToSend)
}

func (h *HttpApiHandler) Login(jsonToSend GenericJson) (result HttpApiResult, err error) {
	url := fmt.Sprintf("%s/session/login", h.Config.Server)
	return h.makeRequest(url, http.MethodPost, jsonToSend)
}

func (h *HttpApiHandler) SetForContainer(jsonToSend GenericJson, containerName, nameSpace string) (result HttpApiResult, err error) {
	url := fmt.Sprintf("%s/namespaces/%s/container-setimage/%s", h.Config.Server, nameSpace, containerName)
	return h.makeRequest(url, http.MethodPatch, jsonToSend)
}

func (h *HttpApiHandler) SetForDeploy(jsonToSend GenericJson, deploy, nameSpace string) (result HttpApiResult, err error) {
	url := fmt.Sprintf("%s/namespaces/%s/deployments/%s/spec", h.Config.Server, nameSpace, deploy)
	return h.makeRequest(url, http.MethodPatch, jsonToSend)
}

func (h *HttpApiHandler) Replace(jsonToSend GenericJson, nameSpace, kind string) (result HttpApiResult, err error) {
	kindI, hasKind := jsonToSend["kind"]
	if !hasKind {
		return result, fmt.Errorf("replace: kind not specified")
	}
	kind, ok := kindI.(string)
	if !ok {
		return result, fmt.Errorf("replace: kind is not a string")
	}
	metadataI, hasMd := jsonToSend["metadata"]
	if !hasMd {
		return result, fmt.Errorf("replace: metadata not specified")
	}
	metadata, ok := metadataI.(map[string]interface{})
	if !ok {
		return result, fmt.Errorf("replace: metadata is not object")
	}
	nameI, hasName := metadata["name"]
	if !hasName {
		return result, fmt.Errorf("replace: metadata has no name attrubute")
	}
	name, ok := nameI.(string)
	if !ok {
		return result, fmt.Errorf("repace: name is not a string")
	}
	url := fmt.Sprintf("%s/namespaces/%s/%s/%s", h.Config.Server, nameSpace, kind, name)
	return h.makeRequest(url, http.MethodPut, jsonToSend)
}

func (h *HttpApiHandler) ReplaceNameSpaces(jsonToSend GenericJson) (result HttpApiResult, err error) {
	metadataI, hasMd := jsonToSend["metadata"]
	if !hasMd {
		return result, fmt.Errorf("replaceNameSpace: metadata not specified")
	}
	metadata, ok := metadataI.(map[string]interface{})
	if !ok {
		return result, fmt.Errorf("replaceNameSpace: metadata is not object")
	}
	nameI, hasName := metadata["name"]
	if !hasName {
		return result, fmt.Errorf("replaceNameSpace: metadata has no name attrubute")
	}
	name, ok := nameI.(string)
	if !ok {
		return result, fmt.Errorf("repaceNameSpace: name is not a string")
	}

	url := fmt.Sprintf("%s/namespaces/%s", h.Config.Server, name)
	return h.makeRequest(url, http.MethodPut, jsonToSend)
}

func (h *HttpApiHandler) Run(jsonToSend GenericJson, nameSpace string) (result HttpApiResult, err error) {
	url := fmt.Sprintf("%s/namespaces/%s/deployments", h.Config.Server, nameSpace)
	return h.makeRequest(url, http.MethodPost, jsonToSend)
}

func (h *HttpApiHandler) Delete(kind, name, nameSpace string, allPods bool) (result HttpApiResult, err error) {
	var url string
	if kind == KindDeployments && allPods {
		url = fmt.Sprintf("%s/namespaces/%s/%s/%s/pods", h.Config.Server, nameSpace, kind, name)
	} else {
		url = fmt.Sprintf("%s/namespaces/%s/%s/%s", h.Config.Server, nameSpace, kind, name)
	}
	return h.makeRequest(url, http.MethodDelete, nil)
}

func (h *HttpApiHandler) DeleteNameSpaces(name string) (result HttpApiResult, err error) {
	url := fmt.Sprintf("%s/namespaces/%s", h.Config.Server, name)
	result, err = h.makeRequest(url, http.MethodDelete, nil)
	return
}

func (h *HttpApiHandler) Get(kind, name, nameSpace string) (result HttpApiResult, err error) {
	var url string
	if name != "" {
		url = fmt.Sprintf("%s/namespaces/%s/%s/%s", h.Config.Server, nameSpace, kind, name)
	} else {
		url = fmt.Sprintf("%s/namespaces/%s/%s", h.Config.Server, nameSpace, kind)
	}
	return h.makeRequest(url, http.MethodGet, nil)
}

func (h *HttpApiHandler) GetNameSpaces(name string) (result HttpApiResult, err error) {
	var url string
	if name != "" {
		url = fmt.Sprintf("%s/namespaces/%s", h.Config.Server, name)
	} else {
		url = fmt.Sprintf("%s/namespaces", h.Config.Server)
	}
	return h.makeRequest(url, http.MethodGet, nil)
}

func (h *HttpApiHandler) GetVolume(name string) (result interface{}, err error) {
	var url string
	if name != "" {
		url = fmt.Sprintf("%s/volumes/%s", h.Config.Server, name)
	} else {
		url = fmt.Sprintf("%s/volumes", h.Config.Server)
	}
	err = h.makeRequestGeneric(url, http.MethodGet, nil, &result)
	return
}

func (apiResult *HttpApiResult) HandleApiResult() error {
	if apiResult == nil {
		return nil
	}
	errCont, hasErr := (*apiResult)["error"]
	if hasErr {
		return fmt.Errorf("api error: %v", errCont)
	}
	return nil
}
