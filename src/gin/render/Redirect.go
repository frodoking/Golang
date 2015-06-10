package render

import (
	"fmt"
	"net/http"
)

type Redirect struct {
	Code     int
	Request  *http.Request
	Location string
}

func (this *Redirect) Write(w http.ResponseWriter) error {
	if this.Code < 300 || this.Code > 308 {
		panic(fmt.Sprintf("Cannot redirect with status code %d", this.Code))
	}

	http.Redirect(w, this.Request, this.Location, this.Code)
	return nil
}
