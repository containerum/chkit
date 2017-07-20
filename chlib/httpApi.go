package chlib

import (
	"bytes"
	"chkit-v2/chlib/dbconfig"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	jww "github.com/spf13/jwalterweatherman"
)

type HttpApiHandler struct {
	cfg     dbconfig.HttpApiConfig
	headers map[string]string
	np      *jww.Notepad
}

type HttpApiResult map[string]interface{}

func NewHttpApiHandler(cfg dbconfig.HttpApiConfig, uuid, token string, np *jww.Notepad) *HttpApiHandler {
	handler := HttpApiHandler{
		cfg: cfg,
		np:  np,
	}
	handler.np.SetPrefix("HTTP")
	handler.headers = make(map[string]string)
	handler.headers["Channel"] = uuid
	handler.headers["Authorization"] = token
	return &handler
}

func (h *HttpApiHandler) makeRequest(url, method string, jsonToSend GenericJson) (result HttpApiResult, err error) {
	client := http.Client{Timeout: h.cfg.Timeout}
	marshalled, _ := json.Marshal(jsonToSend)
	h.np.DEBUG.Printf("%s %s\n", method, url)
	request, err := http.NewRequest(method, url, bytes.NewBuffer(marshalled))
	for k, v := range h.headers {
		request.Header.Add(k, v)
	}
	if err != nil {
		return result, fmt.Errorf("http request create error: %s", err)
	}
	resp, err := client.Do(request)
	if err != nil {
		return result, fmt.Errorf("http request execute error: %s", err)
	}
	h.np.DEBUG.Println("Result", resp.Status)
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}

func (h *HttpApiHandler) Create(jsonToSend GenericJson, kind, nameSpace string) (result HttpApiResult, err error) {
	kind = fmt.Sprintf("%ss", strings.ToLower(kind))
	url := fmt.Sprintf("%s/namespaces/%s/%s", h.cfg.Server, nameSpace, kind)
	return h.makeRequest(url, http.MethodPost, jsonToSend)
}

func (h *HttpApiHandler) Expose(jsonToSend GenericJson, nameSpace string) (result HttpApiResult, err error) {
	url := fmt.Sprintf("%s/namespaces/%s/services", h.cfg.Server, nameSpace)
	return h.makeRequest(url, http.MethodPost, jsonToSend)
}

func (h *HttpApiHandler) Login(jsonToSend GenericJson) (result HttpApiResult, err error) {
	url := fmt.Sprintf("%s/session/login", h.cfg.Server)
	result, err = h.makeRequest(url, http.MethodPost, jsonToSend)
	return
}

func (h *HttpApiHandler) Set(jsonToSend GenericJson, name, nameSpace string) (result HttpApiResult, err error) {
	var url string
	if _, hasReplicas := jsonToSend["replicas"]; hasReplicas {
		url = fmt.Sprintf("%s/namespaces/%s/deployments/%s/spec", h.cfg.Server, nameSpace, name)
	} else {
		url = fmt.Sprintf("%s/namespaces/%s/container/%s", h.cfg.Server, nameSpace, name)
	}
	result, err = h.makeRequest(url, http.MethodPatch, jsonToSend)
	return
}

func (h *HttpApiHandler) Scale(jsonToSend GenericJson, name, nameSpace string) (result HttpApiResult, err error) {
	url := fmt.Sprintf("%s/namespaces/%s/deployments/%s/spec", h.cfg.Server, nameSpace, name)
	result, err = h.makeRequest(url, http.MethodPatch, jsonToSend)
	return
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
	url := fmt.Sprintf("%s/namespaces/%s/%s/%s", h.cfg.Server, nameSpace, kind, name)
	result, err = h.makeRequest(url, http.MethodPut, jsonToSend)
	return
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

	url := fmt.Sprintf("%s/namespaces/%s", h.cfg.Server, name)
	result, err = h.makeRequest(url, http.MethodPut, jsonToSend)
	return
}

func (h *HttpApiHandler) Run(jsonToSend GenericJson, nameSpace string) (result HttpApiResult, err error) {
	url := fmt.Sprintf("%s/namespaces/%s/deployments", h.cfg.Server, nameSpace)
	result, err = h.makeRequest(url, http.MethodPost, jsonToSend)
	return
}

func (h *HttpApiHandler) Delete(kind, name, nameSpace string, allPods bool) (result HttpApiResult, err error) {
	var url string
	if kind == KindDeployments && allPods {
		url = fmt.Sprintf("%s/namespaces/%s/%s/%s/pods", h.cfg.Server, nameSpace, kind, name)
	} else {
		url = fmt.Sprintf("%s/namespaces/%s/%s/%s", h.cfg.Server, nameSpace, kind, name)
	}
	result, err = h.makeRequest(url, http.MethodDelete, nil)
	return
}

func (h *HttpApiHandler) DeleteNameSpaces(name string) (result HttpApiResult, err error) {
	url := fmt.Sprintf("%s/namespaces/%s", h.cfg.Server, name)
	result, err = h.makeRequest(url, http.MethodDelete, nil)
	return
}

func (h *HttpApiHandler) Get(kind, name, nameSpace string) (result HttpApiResult, err error) {
	var url string
	if name != "" {
		url = fmt.Sprintf("%s/namespaces/%s/%s/%s", h.cfg.Server, nameSpace, kind, name)
	} else {
		url = fmt.Sprintf("%s/namespaces/%s/%s", h.cfg.Server, nameSpace, kind)
	}
	result, err = h.makeRequest(url, http.MethodGet, nil)
	return
}

func (h *HttpApiHandler) GetNameSpaces(name string) (result HttpApiResult, err error) {
	var url string
	if name != "" {
		url = fmt.Sprintf("%s/namespaces/%s", h.cfg.Server, name)
	} else {
		url = fmt.Sprintf("%s/namespaces", h.cfg.Server)
	}
	result, err = h.makeRequest(url, http.MethodGet, nil)
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
