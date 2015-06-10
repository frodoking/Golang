package binding

import (
	"net/http"
)

type FormBinding struct{}

func (this *FormBinding) Name() string {
	return "form"
}

func (this *FormBinding) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	req.ParseMultipartForm(32 << 10) //32M
	if err := mapForm(obj, req.Form); err != nil {
		return err
	}

	return validate(obj)
}
