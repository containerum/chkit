package chkitErrors

import (
	bytes "bytes"
	cherry "git.containerum.net/ch/kube-client/pkg/cherry"
	template "text/template"
)

const ()

func ErrUnableToInitClient(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "unable to init client", StatusHTTP: 418, ID: cherry.ErrID{SID: 0x309, Kind: 0x1}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrInvalidUsername(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "invalid username", StatusHTTP: 418, ID: cherry.ErrID{SID: 0x309, Kind: 0x2}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrInvalidPassword(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "invalid password", StatusHTTP: 418, ID: cherry.ErrID{SID: 0x309, Kind: 0x3}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableToReadUsername(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "unable to read username", StatusHTTP: 418, ID: cherry.ErrID{SID: 0x309, Kind: 0x4}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableToReadPassword(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "unable to read password", StatusHTTP: 418, ID: cherry.ErrID{SID: 0x309, Kind: 0x5}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableToSaveLogin(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "unable to save login", StatusHTTP: 418, ID: cherry.ErrID{SID: 0x309, Kind: 0x6}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}
func renderTemplate(templText string) string {
	buf := &bytes.Buffer{}
	templ, err := template.New("").Parse(templText)
	if err != nil {
		return err.Error()
	}
	err = templ.Execute(buf, map[string]string{})
	if err != nil {
		return err.Error()
	}
	return buf.String()
}
