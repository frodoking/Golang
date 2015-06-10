package gin

import (
	"io"
	"os"
)

import (
	"github.com/mattn/go-colorable"
)

import (
	"gin/binding"
)

const (
	ENV_GIN_MODE = "GIN_MODE"
)

const (
	DebugMode   string = "debug"
	ReleaseMode string = "release"
	TestMode    string = "test"
)

const (
	debugCode   = iota
	releaseCode = iota
	testCode    = iota
)

var DefaultWriter io.Writer = colorable.NewColorableStdout()
