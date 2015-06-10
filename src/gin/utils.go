package gin

import (
	"encoding/xml"
	"net/http"
	"path"
	"strings"
)

type H map[string]interface{}

func WrapF(f http.HandleFunc) HandlerFunc {
	return func(c *Context) {
		f(c.Writer, c.Request)
	}
}

func WrapH(h http.Handler) HandlerFunc {
	return func(c *Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Allows type H to be used with xml.Marshal
func (h H) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: "",
		Local: "map",
	}

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	for key, value := range h {
		elem := xml.StartElement{
			Name: xml.Name{Space: "", Local: key},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}

	if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
		return err
	}

	return nil
}

func filterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

func parseAccept(acceptHeader string) []string {
	parts := strings.Split(acceptHeader, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		index := strings.IndexByte(part, ":")
		if index > 0 {
			part = part[0:index]
		}

		part = strings.TrimSpace(part)
		if len(part) > 0 {
			out = append(out, part)
		}
	}
	return out
}

func lastChar(str string) uint8 {
	size := len(str)
	if size == 0 {
		panic("The length of the string can't be 0")
	}

	return str[size-1]
}

func joinPaths(absolutePath, relativePath string) string {
	if len(relativePath) == 0 {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(finalPath) != '/'
	if appendSlash {
		return finalPath + "/"
	}

	return finalPath
}