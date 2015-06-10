package render

import (
	"net/http/httptest"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestRenderJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]interface{}{
		"foo": "bar",
	}
	err := (&JSON{data}).Write(w)

	assert.NoError(t, err)
	assert.Equal(t, w.Body.String(), "{\"foo\":\"bar\"}\n")
	assert.Equal(t, w.Header().Get("Content-Type"), "application/json; chartset=utf-8")
}
