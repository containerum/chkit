package chlib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/kfeofantov/chkit-v2/helpers"
)

type HttpApiConfig struct {
	Server  string            `mapconv:"server"`
	headers map[string]string `mapconv:"-"`
	Timeout time.Duration     `mapconv:"timeout"`
	Uuid    string            `mapconv:"-"`
	Token   string            `mapconv:"-"`
}

type HttpApiHandler struct {
	cfg HttpApiConfig
}

type HttpApiResult map[string]interface{}

const httpApiBucket = "httpApi"

func init() {
	cfg := HttpApiConfig{
		Server:  "http://0.0.0.0",
		Timeout: 10 * time.Second,
	}
	initializers[httpApiBucket] = helpers.StructToMap(cfg)
}

func GetHttpApiCfg() (cfg HttpApiConfig, err error) {
	m, err := readFromBucket(httpApiBucket)
	if err != nil {
		return cfg, fmt.Errorf("load http api config: %s", err)
	}
	err = helpers.FillStruct(&cfg, m)
	if err != nil {
		return cfg, fmt.Errorf("http api config fill: %s", err)
	}
	return cfg, nil
}

func UpdateHttpApiCfg(cfg HttpApiConfig) error {
	return pushToBucket(httpApiBucket, helpers.StructToMap(cfg))
}

func NewHttpApiHandler(cfg HttpApiConfig) *HttpApiHandler {
	handler := HttpApiHandler{cfg: cfg}
	handler.cfg.headers = make(map[string]string)
	handler.cfg.headers["Channel"] = cfg.Uuid
	handler.cfg.headers["Authorization"] = cfg.Token
	return &handler
}

func (h *HttpApiHandler) makeRequest(url, method string, jsonToSend GenericJson) (result HttpApiResult, err error) {
	client := http.Client{Timeout: h.cfg.Timeout}
	marshalled, _ := json.Marshal(jsonToSend)
	request, err := http.NewRequest(method, url, bytes.NewBuffer(marshalled))
	for k, v := range h.cfg.headers {
		request.Header.Add(k, v)
	}
	if err != nil {
		return result, fmt.Errorf("http request create error: %s", err)
	}
	resp, err := client.Do(request)
	if err != nil {
		return result, fmt.Errorf("http request execute error: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return result, err
		} else {
			return result, fmt.Errorf("%s", data)
		}
	}
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
	url := fmt.Sprintf("%s/namespaces/%s/deployment", h.cfg.Server, nameSpace)
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
