package chlib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kfeofantov/chkit-v2/helpers"
)

type HttpApiConfig struct {
	Server  string            `mapconv:"server"`
	Headers map[string]string `mapconv:"headers"`
	Timeout time.Duration     `mapconv:"timeout"`
	Uuid    string            `mapconv:"-"`
}

type HttpApiHandler struct {
	cfg HttpApiConfig
}

type HttpApiResult map[string]interface{}

const httpApiBucket = "httpApi"

func init() {
	cfg := HttpApiConfig{Headers: map[string]string{"Authorization": ""}}
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
	handler.cfg.Headers["Channel"] = cfg.Uuid
	return &handler
}

func (h *HttpApiHandler) makeRequest(url, method string, jsonToSend json.RawMessage) (result HttpApiResult, err error) {
	client := http.Client{Timeout: h.cfg.Timeout}
	request, err := http.NewRequest(method, url, bytes.NewBuffer(jsonToSend))
	for k, v := range h.cfg.Headers {
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
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}

func (h *HttpApiHandler) Login(jsonToSend json.RawMessage) (result HttpApiResult, err error) {
	url := fmt.Sprintf("%s/session/login", h.cfg.Server)
	result, err = h.makeRequest(url, http.MethodPost, jsonToSend)
	return
}
