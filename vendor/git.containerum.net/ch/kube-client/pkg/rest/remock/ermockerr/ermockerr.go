package ermockerr

import (
	bytes "bytes"
	cherry "git.containerum.net/ch/kube-client/pkg/cherry"
	template "text/template"
)

const ()

func ErrInvalidDeployment(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "invalid deployment", StatusHTTP: 418, ID: cherry.ErrID{SID: 0x378, Kind: 0x1}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrInvalidContainer(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "invalid container", StatusHTTP: 418, ID: cherry.ErrID{SID: 0x378, Kind: 0x2}, Details: []string(nil)}
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
