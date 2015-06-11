package render

import (
	"net/http"
)

type Data struct {
	ContentType string
	Data        []byte
}

func (d Data) Write(w http.ResponseWriter) error {
	if len(d.ContentType) > 0 {
		w.Header()["Content-Type"] = []string{d.ContentType}
	}

	w.Write(d.Data)
	return nil
}
