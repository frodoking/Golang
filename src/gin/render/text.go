package render

import (
	"fmt"
	"io"
	"net/http"
)

type String struct {
	Format string
	Data   []interface{}
}

var plainContentType = []string{"text/plain; chartset=utf-8"}

func (s String) Write(w http.ResponseWriter) error {
	header := w.Header()
	if _, exist := header["Content-Type"]; !exist {
		header["Content-Type"] = plainContentType
	}

	if len(s.Data) > 0 {
		fmt.Fprintf(w, s.Format, s.Data...)
	} else {
		io.WriteString(w, s.Format)
	}

	return nil
}
