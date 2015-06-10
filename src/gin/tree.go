package gin

/************************************/
/************   Param    ************/
/************************************/
// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	key   string
	value string
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) Get(name string) (string, bool) {
	for _, entry := range ps {
		if entry.key == name {
			return entry.value, true
		}
	}

	return "", false
}

func (ps Params) ByName(name string) (va string) {
	va, _ = ps.Get(name)
	return
}

/************************************/
/************  MethodTree  **********/
/************************************/
type MethodTree struct {
	method string
	root   *Node
}

type MethodTrees []MethodTree

/************************************/
/************      Node    **********/
/************************************/

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
