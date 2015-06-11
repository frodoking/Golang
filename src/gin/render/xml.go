package render

import (
	"encoding/xml"
	"net/http"
)

type XML struct {
	Data interface{}
}

var xmlContentType = []string{"application/xml; chartset=utf-8"}

func (x XML) Write(w http.ResponseWriter) error {
	w.Header()["Content-Type"] = xmlContentType
	return xml.NewEncoder(w).Encode(x.Data)
}
