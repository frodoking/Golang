package scheduler

import (
	"container/list"
	"crypto/md5"
	"sync"
)

import (
	"go_spider/core/common/request"
)

type QueueScheduler struct {
	locker *sync.Mutex
	rm     bool
	rmKey  map[[md5.Size]byte]*list.Element
	queue  *list.List
}

func NewQueueScheduler(rmDuplicate bool) *QueueScheduler {
	queue := list.New()
	rmKey := make(map[[md5.Size]byte]*list.Element)
	locker := new(sync.Mutex)
	return &QueueScheduler{rm: rmDuplicate, queue: queue, rmKey: rmKey, locker: locker}
}

func (this *QueueScheduler) Push(req *request.Request) {
	this.locker.Lock()
	var key [md5.Size]byte
	if this.rm {
		key = md5.Sum([]byte(req.GetUrl()))
		if _, ok := this.rmKey[key]; ok {
			this.locker.Unlock()
			return
		}
	}

	e := this.queue.PushBack(req)
	if this.rm {
		this.rmKey[key] = e
	}
	this.locker.Unlock()
}
func (this *QueueScheduler) Poll() *request.Request {
	this.locker.Lock()
	if this.queue.Len() <= 0 {
		this.locker.Unlock()
		return nil
	}

	e := this.queue.Front()
	req := e.Value.(*request.Request)
	key := md5.Sum([]byte(req.GetUrl()))
	this.queue.Remove(e)
	if this.rm {
		delete(this.rmKey, key)
	}
	this.locker.Unlock()
	return req
}
func (this *QueueScheduler) Count() int {
	this.locker.Lock()
	len := this.queue.Len()
	this.locker.Unlock()
	return len
}
