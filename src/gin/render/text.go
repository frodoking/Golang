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

func (this *String) Write(w http.ResponseWriter) error {
	header := w.Header()
	if _, exist := header["Content-Type"]; !exist {
		header["Content-Type"] = plainContentType
	}

	if len(this.Data) > 0 {
		fmt.Fprintf(w, this.Format, this.Data...)
	} else {
		io.WriteString(w, this.Format)
	}

	return nil
}
