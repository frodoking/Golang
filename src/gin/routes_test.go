package gin

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func init() {
	SetMode(DebugMode)
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func testRouteOK(method string, t *testing.T) {
	passed := false
	passedAny := false
	r := New()
	r.Any("test2", func(c *Context) {
		passedAny = true
	})

	r.Handle(method, "test", func(c *Context) {
		passed = true
	})

	w := performRequest(r, method, "/test")
	assert.True(t, passed)
	assert.Equal(t, w.Code, http.StatusOK)

	performRequest(r, method, "/test2")
	assert.True(t, passedAny)
}

// TestSingleRouteOK tests that POST route is correctly invoked.
func testRouteNotOK(method string, t *testing.T) {
	passed := false
	router := New()
	router.Handle(method, "/test_2", func(c *Context) {
		passed = true
	})

	w := performRequest(router, method, "/test")

	assert.False(t, passed)
	assert.Equal(t, w.Code, http.StatusNotFound)
}

// TestSingleRouteOK tests that POST route is correctly invoked.
func testRouteNotOK2(method string, t *testing.T) {
	passed := false
	router := New()
	router.HandleMethodNotAllowed = true
	var methodRoute string
	if method == "POST" {
		methodRoute = "GET"
	} else {
		methodRoute = "POST"
	}
	router.Handle(methodRoute, "/test", func(c *Context) {
		passed = true
	})

	w := performRequest(router, method, "/test")

	assert.False(t, passed)
	assert.Equal(t, w.Code, http.StatusMethodNotAllowed)
}

func TestRouterGroupRouteOK(t *testing.T) {
	testRouteOK("GET", t)
	testRouteOK("POST", t)
	testRouteOK("PUT", t)
	testRouteOK("PATCH", t)
	testRouteOK("HEAD", t)
	testRouteOK("OPTIONS", t)
	testRouteOK("DELETE", t)
	testRouteOK("CONNECT", t)
	testRouteOK("TRACE", t)
}

// TestSingleRouteOK tests that POST route is correctly invoked.
func TestRouteNotOK(t *testing.T) {
	testRouteNotOK("GET", t)
	testRouteNotOK("POST", t)
	testRouteNotOK("PUT", t)
	testRouteNotOK("PATCH", t)
	testRouteNotOK("HEAD", t)
	testRouteNotOK("OPTIONS", t)
	testRouteNotOK("DELETE", t)
	testRouteNotOK("CONNECT", t)
	testRouteNotOK("TRACE", t)
}

// TestSingleRouteOK tests that POST route is correctly invoked.
func TestRouteNotOK2(t *testing.T) {
	testRouteNotOK2("GET", t)
	testRouteNotOK2("POST", t)
	testRouteNotOK2("PUT", t)
	testRouteNotOK2("PATCH", t)
	testRouteNotOK2("HEAD", t)
	testRouteNotOK2("OPTIONS", t)
	testRouteNotOK2("DELETE", t)
	testRouteNotOK2("CONNECT", t)
	testRouteNotOK2("TRACE", t)
}

// TestContextParamsGet tests that a parameter can be parsed from the URL.
func TestRouteParamsByName(t *testing.T) {
	name := ""
	lastName := ""
	wild := ""
	router := New()
	router.GET("/test/:name/:last_name/*wild", func(c *Context) {
		name = c.Params.ByName("name")
		lastName = c.Params.ByName("last_name")

		var ok bool
		wild, ok = c.Params.Get("wild")

		assert.True(t, ok)
		assert.Equal(t, name, c.Param("name"))
		assert.Equal(t, lastName, c.Param("last_name"))

		assert.Empty(t, c.Param("xxx"))
		assert.Empty(t, c.Params.ByName("xxx"))

		sss, ok := c.Params.Get("sss")
		assert.Empty(t, sss)
		assert.False(t, ok)
	})

	w := performRequest(router, "GET", "/test/john/smith/is/super/great")

	assert.Equal(t, w.Code, 200)
	assert.Equal(t, name, "john")
	assert.Equal(t, lastName, "smith")
	assert.Equal(t, wild, "is/super/great")
}

// TestHandleStaticFile - ensure the static file handles properly
func TestRouteStaticFile(t *testing.T) {
	// SETUP file
	testRoot, _ := os.Getwd()
	f, err := ioutil.TempFile(testRoot, "")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f.Name())
	f.WriteString("Gin Web Framework")
	f.Close()

	dir, filename := filepath.Split(f.Name())

	// SETUP gin
	router := New()
	router.Static("/using_static", dir)
	router.StaticFile("/result", f.Name())

	w := performRequest(router, "GET", "/using_static/"+filename)
	w2 := performRequest(router, "GET", "/result")
	debugPrint("current file:%s", f)
	debugPrint("current w:%s", w)

	assert.Equal(t, w, w2)
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "Gin Web Framework")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/plain; charset=utf-8")

	w3 := performRequest(router, "HEAD", "/using_static/"+filename)
	w4 := performRequest(router, "HEAD", "/result")

	assert.Equal(t, w3, w4)
	assert.Equal(t, w3.Code, 200)
}
