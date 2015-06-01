package scheduler

import (
	"go_spider/core/common/request"
)

type Scheduler interface {
	Push(req *request.Request)
	Poll() *request.Request
	Count() int
}
