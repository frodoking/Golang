package binding

import (
	"net/http"
)

const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
)

type Binding interface {
	Name() string
	Bind(*http.Request, interface{}) error
}

type StructValidator interface {
	ValidateStruct(interface{}) error
}

var Validator StructValidator = &DefaultValidator{}

var (
	JSON = &JSONBinding{}
	XML  = &XMLBinding{}
	Form = &FormBinding{}
)

func Default(method, contentType string) Binding {
	if method == "GET" {
		return Form
	} else {
		switch contentType {
		case MIMEJSON:
			return JSON
		case MIMEXML, MIMEXML2:
			return XML
			//case MIMEPOSTForm, MIMEMultipartPOSTForm:
		default:
			return Form
		}
	}
}

func validate(obj interface{}) error {
	if Validator == nil {
		return nil
	}

	return Validator.ValidateStruct(obj)
}
