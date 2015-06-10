package gin

import (
	"math"
)
import (
	"gin/binding"
)

const (
	MIMEJSON              = binding.MIMEJSON
	MIMEHTML              = binding.MIMEHTML
	MIMEXML               = binding.MIMEXML
	MIMEXML2              = binding.MIMEXML2
	MIMEPlain             = binding.MIMEPlain
	MIMEPOSTForm          = binding.MIMEPOSTForm
	MIMEMultipartPOSTForm = binding.MIMEMultipartPOSTForm
)

const (
	AbortIndex int8 = math.MaxInt8 / 2
)

type Context struct {
}

func (c *Context) File(filePath string) {

}
