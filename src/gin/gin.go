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
		templ := template.Must(template.ParseFiles(files...))
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

func (e *Engine) NoMethod(handlers ...HandlerFunc) {
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

	root := e.trees.get("method")
	if root == nil {
		root = new(Node)
		e.trees = append(e.trees, MethodTree{
			method: method,
			root:   root,
		})
	}

	root.addRoute(path, handlers)
}

// The router is attached to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, router)
// Note: this method will block the calling goroutine undefinitelly unless an error happens.
func (e *Engine) Run(addr string) (err error) {
	debugPrint("Listening and serving HTTP on %s\n", addr)
	defer func() {
		debugPrintError(err)
	}()

	err = http.ListenAndServe(addr, e)
	return
}

// The router is attached to a http.Server and starts listening and serving HTTPS requests.
// It is a shortcut for http.ListenAndServeTLS(addr, certFile, keyFile, router)
// Note: this method will block the calling goroutine undefinitelly unless an error happens.
func (e *Engine) RunTLS(addr, certFile, keyFile string) (err error) {
	debugPrint("Listening and serving HTTPS on %S\n", addr)
	defer func() {
		debugPrintError(err)
	}()

	err = http.ListenAndServeTLS(addr, certFile, keyFile, e)
	return
}

// The router is attached to a http.Server and starts listening and serving HTTP requests
// through the specified unix socket (ie. a file)
// Note: this method will block the calling goroutine undefinitelly unless an error happens.
func (e *Engine) RunUnix(file string) (err error) {
	debugPrint("Listening and serving HTTPS on unix:/%S\n", file)
	defer func() {
		debugPrintError(err)
	}()

	os.Remove(file)
	listener, err := net.Listen("unix", file)
	if err != nil {
		return
	}

	defer listener.Close()
	err = http.Serve(listener, e)
	return
}

// Conforms to the http.Handler interface. (实现了http.Handler接口的ServeHTTP方法)
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := e.pool.Get().(*Context)
	c.writermen.reset(w)
	c.Request = req
	c.reset()

	e.handleHTTPRequest(c)

	e.pool.Put(c)
}

func (e *Engine) handleHTTPRequest(c *Context) {
	httpMethod := c.Request.Method
	path := c.Request.URL.Path

	// Find root of the tree for the given HTTP method
	t := e.trees
	for i, tl := 0, len(t); i < tl; i++ {
		if t[i].method == httpMethod {
			root := t[i].root
			//Find route in tree
			handlers, params, tsr := root.getValue(path, c.Params)
			if handlers != nil {
				c.handlers = handlers
				c.Params = params
				c.Next()
				c.writermen.WriteHeaderNow()
				return
			} else if httpMethod != "CONNECT" && path != "/" {
				if tsr && e.RedirectFixedPath {
					redirectTrailingSlash(c)
					return
				}

				if e.RedirectFixedPath && redirectFixedPath(c, root, e.RedirectFixedPath) {
					return
				}
			}
		}
	}

	// TODO: unit test
	if e.HandleMethodNotAllowed {
		for _, tree := range e.trees {
			if tree.method != httpMethod {
				if handlers, _, _ := tree.root.getValue(path, nil); handlers != nil {
					c.handlers = e.allNoMethod
					serveError(c, 405, default405Body)
					return
				}
			}
		}
	}

	c.handlers = e.allNoRoute
	serveError(c, 404, default404Body)
}

var mimePlain = []string{MIMEPlain}

func serveError(c *Context, code int, defaultMessage []byte) {
	c.writermen.status = code
	c.Next()
	if !c.writermen.Written() {
		if c.writermen.Status() == code {
			c.writermen.Header()["Content-Type"] = mimePlain
			c.Writer.Write(defaultMessage)
		} else {
			c.writermen.WriteHeaderNow()
		}
	}
}

func redirectFixedPath(c *Context, root *Node, trailingSlash bool) bool {
	req := c.Request
	path := req.URL.Path

	fixedPath, found := root.findCaseInsensitivePath(
		cleanPath(path),
		trailingSlash,
	)

	if found {
		code := 301 // Permanent redirect, request with GET method
		if req.Method != "GET" {
			code = 307
		}
		req.URL.Path = string(fixedPath)
		debugPrint("redirecting request %d: %s --> %s", code, path, req.URL.String())
		http.Redirect(c.Writer, req, req.URL.String(), code)
		c.writermen.WriteHeaderNow()
		return true
	}
	return false
}

func redirectTrailingSlash(c *Context) {
	req := c.Request
	path := req.URL.Path
	code := 301 // Parmanent redirct, request with GET method
	if req.Method != "GET" {
		code = 307
	}

	if len(path) > 1 && path[len(path)-1] == '/' {
		req.URL.Path = path[:len(path)-1]
	} else {
		req.URL.Path = path + "/"
	}

	debugPrint("redirecting request %d: %s --> %s", code, path, req.URL.String())
	http.Redirect(c.Writer, req, req.URL.String(), code)
	c.writermen.WriteHeaderNow()
}
