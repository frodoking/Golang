package mlog

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

type filelog struct {
	plog
	loginst *log.Logger
}

var flog *filelog

// 单例
func LogInst() *filelog {
	if flog == nil {
		InitFilelog(false, "")
	}
	return flog
}

func InitFilelog(isopen bool, fp string) {
	if !isopen {
		flog = &filelog{}
		flog.loginst = nil
		flog.isopen = isopen
		return
	}

	if fp == "" {
		wd := os.Getenv("GOPATH")
		if wd == "" {
			panic("GOPATH is not setted in env.")
			file, _ := exec.LookPath(os.Args[0])
			path := filepath.Dir(file)
			wd = path
		}

		if wd == "" {
			panic("GOPATH is not setted in env or can not get exe path.")
		}

		fp = wd + "/log/"
	}
	flog = newFilelog(isopen, fp)
}

// The newFilelog returns initialized filelog object.
// The default file path is "WORKDIR/log/log.2011-01-01".
func newFilelog(isopen bool, logpath string) *filelog {
	year, month, day := time.Now().Date()
	filename := "log." + strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)
	err := os.MkdirAll(logpath, 0755)
	if err != nil {
		panic("logpath error : " + logpath + "\n")
	}

	f, err := os.OpenFile(logpath+"/"+filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic("log file open error : " + logpath + "/" + filename + "\n")
	}

	pfilelog := &filelog{}
	pfilelog.loginst = log.New(f, "", log.LstdFlags)
	pfilelog.isopen = isopen
	return pfilelog
}

func (this *filelog) LogError(str string) {
	this.log("[ERROR]", str)
}

func (this *filelog) LogInfo(str string) {
	this.log("[INFO]", str)
}

func (this *filelog) log(label string, str string) {
	if !this.isopen {
		return
	}

	file, line := this.getCaller()
	this.loginst.Printf("%s:%d: %s %s\n", file, line, label, str)
}
