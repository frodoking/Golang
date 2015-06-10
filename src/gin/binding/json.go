package binding

import (
	"encoding/json"
	"net/http"
)

type JSONBinding struct{}

func (this *JSONBinding) Name() string {
	return "json"
}

func (this *JSONBinding) Bind(req *http.Request, obj interface{}) error {
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(obj); err != nil {
		return err
	}

	return validate(obj)
}
