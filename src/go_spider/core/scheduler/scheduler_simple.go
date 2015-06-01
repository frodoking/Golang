package scheduler

import (
	"go_spider/core/common/request"
)

type SimpleScheduler struct {
	// 通道
	queue chan *request.Request
}

func NewSimpleScheduler() *SimpleScheduler {
	ch := make(chan *request.Request, 1024)
	return &SimpleScheduler{ch}
}

func (this *SimpleScheduler) Push(req *request.Request) {
	this.queue <- req
}

func (this *SimpleScheduler) Poll() *request.Request {
	if len(this.queue) == 0 {
		return nil
	} else {
		return <-this.queue
	}
}

func (this *SimpleScheduler) Count() int {
	return len(this.queue)
}
