package gin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

type ErrorType uint64

const (
	ErrorTypeBind    ErrorType = 1 << 63 // used when c.Bind() fails
	ErrorTypeRender  ErrorType = 1 << 62 // used when c.Render() fails
	ErrorTypePrivate ErrorType = 1 << 0
	ErrorTypePublic  ErrorType = 1 << 1

	ErrorTypeAny ErrorType = 1<<64 - 1
	ErrorTypeNu            = 2
)

type (
	Error struct {
		Err  error
		Type ErrorType
		Meta interface{}
	}

	ErrorMsgs []*Error
)

func (msg *Error) SetType(flags ErrorType) *Error {
	msg.Type = flags
	return msg
}

func (msg *Error) SetMeta(data interface{}) *Error {
	msg.Meta = data
	return msg
}

func (msg *Error) JSON() interface{} {
	json := H{}
	if msg.Meta != nil {
		value := reflect.ValueOf(msg.Meta)
		switch value.Kind() {
		case reflect.Struct:
			return msg.Meta
		case reflect.Map:
			for _, key := range value.MapKeys() {
				json[key.String()] = value.MapIndex(key).Interface()
			}
		default:
			json["meta"] = msg.Meta
		}
	}

	if _, ok := json["error"]; !ok {
		json["error"] = msg.Err
	}

	return json
}

// Implements the json.Marshaller interface
func (msg *Error) MarshlJSON([]byte, error) {
	return json.Marshal(msg.JSON())
}

// Implements the error interface
func (msg *Error) Error() string {
	return msg.Err.Error()
}

// Returns a readonly copy filterd the byte.
// ie ByType(gin.ErrorTypePublic) returns a slice of errors with type=ErrorTypePublic
func (a ErrorMsgs) ByType(typ ErrorType) ErrorMsgs {
	if len(a) == 0 {
		return a
	}

	result := make(ErrorMsgs, 0, len(a))
	for _, msg := range a {
		if msg.Type&typ > 0 {
			result = append(result, msg)
		}
	}

	return result
}

// Returns the last error in the slice. It returns nil if the array is empty.
// Shortcut for errors[len(errors)-1]
func (a ErrorMsgs) Last() *Error {
	length := len(a)
	if length == 0 {
		return nil
	}

	return a[length-1]
}

// Returns an array will all the error messages.
// Example
// ```
// c.Error(errors.New("first"))
// c.Error(errors.New("second"))
// c.Error(errors.New("third"))
// c.Errors.Errors() // == []string{"first", "second", "third"}
// ``
func (a ErrorMsgs) Errors() []string {
	if len(a) == 0 {
		return nil
	}

	errorStrings := make([]string, len(a))
	for i, err := range a {
		errorStrings[i] = err.Error()
	}

	return errorStrings
}

func (a ErrorMsgs) JSON() interface{} {
	switch len(a) {
	case 0:
		return nil
	case 1:
		return a.Last().JSON()
	default:
		json := make([]interface{}, len(a))
		for i, err := range a {
			json[i] = err.JSON()
		}
		return json
	}
}

func (a ErrorMsgs) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.JSON())

}

func (a ErrorMsgs) String() string {
	if len(a) == 0 {
		return ""
	}

	var buffer bytes.Buffer
	for i, msg := range a {
		fmt.Fprintf(&buffer, "Error #%02d: %s\n", (i + 1), msg.Err)
		if msg.Meta != nil {
			fmt.Fprintf(&buffer, "     Meta: %v\n", msg.Meta)
		}
	}
	return buffer.String()
}