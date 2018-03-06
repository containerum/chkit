package rest

import "strings"

const (
	HeaderUserToken       = "User-Token"
	HeaderUserFingerprint = "User-Client"
)

// P -- URL path params
type P map[string]string

// Q -- URL query params
type Q map[string]string

type URL struct {
	Path   string
	Params P
}

func (u *URL) Build() string {
	addr := u.Path
	for k, v := range u.Params {
		addr = strings.Replace(addr,
			"{"+k+"}", v, -1)
	}
	return addr
}

// Rq -- request params
type Rq struct {
	Result interface{}
	Body   interface{}
	URL    URL
	Query  Q
	Token  string
}

// REST -- rest client interface
type REST interface {
	SetToken(string)
	SetFingerprint(string)
	Get(Rq) error
	Put(Rq) error
	Post(Rq) error
	Delete(Rq) error
}
