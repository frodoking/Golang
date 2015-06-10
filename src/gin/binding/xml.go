package binding

import (
	"encoding/xml"
	"net/http"
)

type XMLBinding struct{}

func (this *XMLBinding) Name() string {
	return "xml"
}

func (this *XMLBinding) Bind(req *http.Request, obj interface{}) error {
	decoder := xml.NewDecoder(req.Body)
	if err := decoder.Decode(obj); err != nil {
		return err
	}

	return validate(obj)
}
