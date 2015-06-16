package gin

import (
	"html/template"
	"net"
	"net/http"
	"os"
	"sync"
)

import (
	"gin/render"
)

const (
	Version = "v.frodo"
)

var default404Body = []byte("404 page not found")
var default405Body = []byte("405 method not allowed")

type (
	HandlerFunc   func(*Context)
	HandlersChain []HandlerFunc

	Engine struct {
		RouterGroup
		HTMLRender  render.HTMLRender
		allNoRoute  HandlersChain
		allNoMethod HandlersChain
		noRoute     HandlersChain
		noMethod    HandlersChain
		pool        sync.Pool
		trees       MethodTrees

		// Enables automatic redirection if the current route can't be matched but a
		// handler for the path with (without) the trailing slash exists.
		// For example if /foo/ is requested but a route only exists for /foo, the
		// client is redirected to /foo with http status code 301 for GET requests
		// and 307 for all other request methods.
		RedirectTrailingSlash bool

		// If enabled, the router tries to fix the current request path, if no
		// handle is registered for it.
		// First superfluous path elements like ../ or // are removed.
		// Afterwards the router does a case-insensitive lookup of the cleaned path.
		// If a handle can be found for this route, the router makes a redirection
		// to the corrected path with status code 301 for GET requests and 307 for
		// all other request methods.
		// For example /FOO and /..//Foo could be redirected to /foo.
		// RedirectTrailingSlash is independent of this option.
		RedirectFixedPath bool

		// If enabled, the router checks if another method is allowed for the
		// current route, if the current request can not be routed.
		// If this is the case, the request is answered with 'Method Not Allowed'
		// and HTTP status code 405.
		// If no other Method is allowed, the request is delegated to the NotFound
		// handler.
		HandleMethodNotAllowed bool
	}
)

// Returns a new blank Engine instance without any middleware attached.
// The most basic configuration
func New() *Engine {
	debugPrintWARNING()
	engine := &Engine{
		RouterGroup: RouterGroup{
			Handlers: nil,
			BasePath: "/",
		},
		trees: make(MethodTrees, 0, 5),
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      false,
		HandleMethodNotAllowed: false,
	}

	engine.RouterGroup.engine = engine
	engine.pool.New = func() interface{} {
		return engine.allocateContext()
	}

	return engine
}

// Returns a Engine instance with the Logger and Recovery already attached.
func Default() *Engine {
	engine := New()
	engine.Use(Recovery(), Logger())
	return engine
}

func (e *Engine) allocateContext() (context *Context) {
	return &Context{engine: e}
}

func (e *Engine) LoadHTMLGlob(pattern string) {
	if IsDebugging() {
		e.HTMLRender = render.HTMLDebug{Glob: pattern}
	} else {
		templ := template.Must(template.ParseGlob(pattern))
		e.SetHTMLTemplate(templ)
	}
}

func (e *Engine) LoadHTMLFiles(files ...string) {
	if IsDebugging() {
		e.HTMLRender = render.HTMLDebug{Files: files}
	} else {
		templ := template.Must(template.ParseFiles(files))
		e.SetHTMLTemplate(templ)
	}
}

func (e *Engine) SetHTMLTemplate(templ *template.Template) {
	e.HTMLRender = render.HTMLProduction{Template: templ}
}

// Adds handlers for NoRoute. It return a 404 code by default.
func (e *Engine) NoRoute(handlers ...HandlerFunc) {
	e.noRoute = handlers
	e.rebuild404Handlers()
}

func (e *Engine) noMethod(handlers ...HandlerFunc) {
	e.noMethod = handlers
	e.rebuild405Handlers()
}

// Attachs a global middleware to the router. ie. the middlewares attached though Use() will be
// included in the handlers chain for every single request. Even 404, 405, static files...
// For example, this is the right place for a logger or error management middleware.
func (e *Engine) Use(middlewares ...HandlerFunc) {
	e.RouterGroup.Use(middlewares...)
	e.rebuild404Handlers()
	e.rebuild405Handlers()
}

func (e *Engine) rebuild404Handlers() {
	e.allNoRoute = e.combineHandlers(e.noRoute)
}

func (e *Engine) rebuild405Handlers() {
	e.allNoMethod = e.combineHandlers(e.noMethod)
}

func (e *Engine) addRoute(method, path string, handlers HandlersChain) {
	debugPrintRoute(method, path, handlers)

	if path[0] != '/' {
		panic("path must begin with '/'")
	}

	if method == "" {
		panic("HTTP method can not be empty")
	}

	if len(handlers) == 0 {
		panic("there must be at least one handler")
	}

	// FIXME
}
