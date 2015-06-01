package page_items

import (
	"go_spider/core/common/request"
)

// PageItems represents an entity save result parsed by PageProcesser and will be output at last.
type PageItems struct {
	req *request.Request
	// The items is the container of parsed result.
	items map[string]string
	// The skip represents whether send ResultItems to scheduler or not.
	skip bool
}

func NewPageItems(req *request.Request) *PageItems {
	items := make(map[string]string)
	return &PageItems{req: req, items: items}
}

func (this *PageItems) GetRequest() *request.Request {
	return this.req
}

func (this *PageItems) AddItem(key string, item string) {
	this.items[key] = item
}

func (this *PageItems) GetItem(key string) (string, bool) {
	t, ok := this.items[key]
	return t, ok
}

func (this *PageItems) GetAll() map[string]string {
	return this.items
}

func (this *PageItems) GetSkip() bool {
	return this.skip
}

func (this *PageItems) SetSkip(skip bool) *PageItems {
	this.skip = skip
	return this
}
