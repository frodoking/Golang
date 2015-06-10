package gin

import (
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
	}
)
