package gin

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func init() {
	SetMode(TestMode)
}

func TestRouterGroupBasic(t *testing.T) {
	router := New()
	group := router.Group("/hola", func(c *Context) {})
	group.Use(func(c *Context) {})

	assert.Len(t, group.Handlers, 2)
	assert.Equal(t, group.BasePath, "/hola")
	assert.Equal(t, group.engine, router)

	group2 := group.Group("manu")
	group2.Use(func(c *Context) {}, func(c *Context) {})

	assert.Len(t, group2.Handlers, 4)
	assert.Equal(t, group2.BasePath, "/hola/manu")
	assert.Equal(t, group2.engine, router)
}

func TestRouterGroupBasicHandle(t *testing.T) {
	performRequestInGroup(t, "GET")
	performRequestInGroup(t, "POST")
	performRequestInGroup(t, "PUT")
	performRequestInGroup(t, "PATCH")
	performRequestInGroup(t, "DELETE")
	performRequestInGroup(t, "HEAD")
	performRequestInGroup(t, "OPTIONS")
}

func performRequestInGroup(t *testing.T, method string) {
	router := New()

	v1 := router.Group("v1", func(c *Context) {})
	assert.Equal(t, v1.BasePath, "/v1")

	login := v1.Group("/login/", func(c *Context) {}, func(c *Context) {})
	assert.Equal(t, login.BasePath, "/v1/login")

	handler := func(c *Context) {
		c.String(400, "the method was %s and index %d", c.Request.Method, c.index)
	}

	switch method {
	case "GET":
		v1.GET("/test", handler)
		login.GET("/test", handler)
	case "POST":
		v1.POST("/test", handler)
		login.POST("/test", handler)
	case "PUT":
		v1.PUT("/test", handler)
		login.PUT("/test", handler)
	case "PATCH":
		v1.PATCH("/test", handler)
		login.PATCH("/test", handler)
	case "DELETE":
		v1.DELETE("/test", handler)
		login.DELETE("/test", handler)
	case "HEAD":
		v1.HEAD("/test", handler)
		login.HEAD("/test", handler)
	case "OPTIONS":
		v1.OPTIONS("/test", handler)
		login.OPTIONS("/test", handler)
	default:
		panic("unknown method")
	}

	w := performRequest(router, method, "/v1/login/test")
	assert.Equal(t, w.Code, 400)
	assert.Equal(t, w.Body.String(), "the method was "+method+" and index 3")

	w = performRequest(router, method, "/v1/test")
	assert.Equal(t, w.Code, 400)
	assert.Equal(t, w.Body.String(), "the method was "+method+" and index 1")
}
