package gin

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

var _ ResponseWriter = &responseWriter{}
var _ http.ResponseWriter = &responseWriter{}
var _ http.ResponseWriter = ResponseWriter(&responseWriter{})
var _ http.Hijacker = ResponseWriter(&responseWriter{})
var _ http.Flusher = ResponseWriter(&responseWriter{})
var _ http.CloseNotifier = ResponseWriter(&responseWriter{})

func init() {
	SetMode(TestMode)
}

func TestResponseWriterReset(t *testing.T) {
	testWritter := httptest.NewRecorder()
	writer := &responseWriter{}
	var w ResponseWriter = writer
	writer.reset(testWritter)
	assert.Equal(t, writer.size, -1)
	assert.Equal(t, writer.status, 200)
	assert.Equal(t, writer.ResponseWriter, testWritter)
	assert.Equal(t, w.Size(), -1)
	assert.Equal(t, w.Status(), 200)
	assert.False(t, w.Written())
}
