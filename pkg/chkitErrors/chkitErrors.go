package chkitErrors

import (
	bytes "bytes"
	template "text/template"

	cherry "git.containerum.net/ch/kube-client/pkg/cherry"
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
