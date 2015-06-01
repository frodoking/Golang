package mlog

import (
	"runtime"
)

type plog struct {
	isopen bool
}

func (*plog) getCaller() (string, int) {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	}
	return file, line
}

func (this *plog) Open() {
	this.isopen = true
}

func (this *plog) Close() {
	this.isopen = false
}
