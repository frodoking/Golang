package render

import (
	"net/http"
)

type Data struct {
	ContentType string
	Data        []byte
}

func (this *Data) Write(w http.ResponseWriter) error {
	if len(this.ContentType) > 0 {
		w.Header()["Content-Type"] = []string{this.ContentType}
	}

	w.Write(this.Data)
	return nil
}
