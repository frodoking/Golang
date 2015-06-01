package mlog

import (
	"log"
	"os"
)

type strace struct {
	plog

	loginst *log.Logger
}

var pstrace *strace

func StraceInst() *strace {
	if pstrace == nil {
		pstrace = newStrace()
	}

	return pstrace
}

func newStrace() *strace {
	pstrace := &strace{}
	pstrace.loginst = log.New(os.Stderr, "", log.LstdFlags)
	pstrace.isopen = true
	return pstrace
}

func (this *strace) Println(str string) {
	if !this.isopen {
		return
	}

	this.loginst.Printf("%s\n", str)
}
