package gin

import (
	"bytes"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

// TestPanicInHandler assert that panic has been recovered.
func TestPanicInHandler(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(RecoveryWithWriter(buffer))
	router.GET("/recovery", func(c *Context) {
		panic("Oupps, Houston, we have a problem")
	})

	//RUN
	w := performRequest(router, "GET", "/recovery")
	//TEST
	assert.Equal(t, w.Code, 500)
	assert.Contains(t, buffer.String(), "Panic recovery -> Oupps, Houston, we have a problem")
	assert.Contains(t, buffer.String(), "TestPanicInHandler")
}
