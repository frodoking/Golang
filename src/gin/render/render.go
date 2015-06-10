package render

import (
	"net/http"
)

type Render interface {
	Write(w http.ResponseWriter) error
}

var (
	_ Render     = &JSON{}
	_ Render     = &IndentedJSON{}
	_ Render     = &XML{}
	_ Render     = &String{}
	_ Render     = &Redirect{}
	_ Render     = &Data{}
	_ Render     = &HTML{}
	_ HTMLRender = &HTMLDebug{}
	_ HTMLRender = &HTMLProduction{}
)
