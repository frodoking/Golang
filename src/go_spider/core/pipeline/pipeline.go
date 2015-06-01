package pipeline

import (
	"go_spider/core/common/com_interfaces"
	"go_spider/core/common/page_items"
)

type Pipeline interface {
	Process(items *page_items.PageItems, t com_interfaces.Task)
}

type CollectionPipline interface {
	Pipeline
	// The GetCollected returns result saved in in process's memory temporarily.
	GetCollected() []*page_items.PageItems
}
