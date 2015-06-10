package gin

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	key   string
	value string
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

type MethodTree struct {
	method string
	root   *Node
}

type MethodTrees []MethodTree

type NodeType uint8

type Node struct {
	path      string
	wildChild bool
	nType     NodeType
	maxParams uint8
	indices   string
	children  []*Node
	handlers  HandlersChain
	priority  uint32
}
