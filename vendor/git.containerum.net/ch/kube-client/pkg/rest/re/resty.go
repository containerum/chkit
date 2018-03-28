package re

import (
	"crypto/tls"
	"io"
	"net/http"

	"git.containerum.net/ch/kube-client/pkg/cherry"
	"git.containerum.net/ch/kube-client/pkg/rest"
	resty "github.com/go-resty/resty"
)

var (
	_ rest.REST = &Resty{}
)

// Resty -- resty client,
// implements REST interface
type Resty struct {
	token, fingerprint string
	client             *resty.Client
}

// NewResty -- Resty constuctor
func NewResty(configs ...func(*Resty)) *Resty {
	re := &Resty{
		client: resty.New().
			SetRESTMode().
			SetHeader("User-Agent", "kube-client"),
	}
	for _, config := range configs {
		config(re)
	}
	return re
}

func WithHost(addr string) func(re *Resty) {
	return func(re *Resty) {
		re.client.SetHostURL(addr)
	}
}

func WithLogger(wr io.Writer) func(re *Resty) {
	return func(re *Resty) {
		re.client.SetDebug(true)
		re.client.SetLogger(wr)
		re.client.Log.Println("rest client in debug mode")
	}
}

func SkipTLSVerify(re *Resty) {
	re.client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
}

// Get -- http get method
func (re *Resty) Get(reqconfig rest.Rq) error {
	resp, err := ToResty(reqconfig, re.client.R()).
		Get(reqconfig.URL.Build())
	if err = rest.MapErrors(resp, err, http.StatusOK); err != nil {
		return err
	}
	if reqconfig.Result != nil {
		rest.CopyInterface(reqconfig.Result, resp.Result())
	}
	return nil
}

// Put -- http put method
func (re *Resty) Put(reqconfig rest.Rq) error {
	resp, err := ToResty(reqconfig, re.client.R()).
		Put(reqconfig.URL.Build())
	if err = rest.MapErrors(resp, err,
		http.StatusOK,
		http.StatusAccepted,
		http.StatusCreated); err != nil {
		return err
	}
	if reqconfig.Result != nil {
		rest.CopyInterface(reqconfig.Result, resp.Result())
	}
	return nil
}

// Post -- http post method
func (re *Resty) Post(reqconfig rest.Rq) error {
	resp, err := ToResty(reqconfig, re.client.R()).
		Post(reqconfig.URL.Build())
	if err = rest.MapErrors(resp, err,
		http.StatusOK,
		http.StatusAccepted,
		http.StatusCreated); err != nil {
		return err
	}
	if reqconfig.Result != nil {
		rest.CopyInterface(reqconfig.Result, resp.Result())
	}
	return nil
}

// Delete -- http delete method
func (re *Resty) Delete(reqconfig rest.Rq) error {
	resp, err := ToResty(reqconfig, re.client.R()).
		Post(reqconfig.URL.Build())
	return rest.MapErrors(resp, err,
		http.StatusOK,
		http.StatusAccepted,
		http.StatusNoContent)
}

func (re *Resty) SetToken(token string) {
	re.client.SetHeader(rest.HeaderUserToken, token)
	re.token = token
}

func (re *Resty) SetFingerprint(fingerprint string) {
	re.client.SetHeader(rest.HeaderUserFingerprint, fingerprint)
	re.fingerprint = fingerprint
}

func (re *Resty) GetToken() string {
	return re.token
}

func (re *Resty) GetFingerprint() string {
	return re.fingerprint
}

// ToResty -- maps Rq data to resty request
func ToResty(rq rest.Rq, req *resty.Request) *resty.Request {
	if rq.Result != nil {
		req = req.SetResult(rq.Result)
	}
	if rq.Body != nil {
		req = req.SetBody(rq.Body)
	}
	if len(rq.Query) > 0 {
		req = req.SetQueryParams(rq.Query)
	}
	return req.SetError(cherry.Err{})
}
