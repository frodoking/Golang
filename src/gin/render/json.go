package render

import (
	"encoding/json"
	"net/http"
)

type (
	JSON struct {
		Data interface{}
	}

	IndentedJSON struct {
		Data interface{}
	}
)

var jsonContentType = []string{"application/json; chartset=utf-8"}

func (this *JSON) Write(w http.ResponseWriter) error {
	w.Header()["Content-Type"] = jsonContentType
	return json.NewEncoder(w).Encode(this.Data)
}

func (this *IndentedJSON) Write(w http.ResponseWriter) error {
	w.Header()["Content-Type"] = jsonContentType
	jsonBytes, err := json.MarshalIndent(this.Data, "", "     ")
	if err != nil {
		return err
	}
	w.Write(jsonBytes)
	return nil
}
