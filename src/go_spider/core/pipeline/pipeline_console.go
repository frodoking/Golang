package pipeline

import (
	"fmt"
)

import (
	"go_spider/core/common/com_interfaces"
	"go_spider/core/common/page_items"
)

type PipelineConsole struct {
}

func NewPipelineConsole() *PipelineConsole {
	return &PipelineConsole{}
}

func (this *PipelineConsole) Process(items *page_items.PageItems, t com_interfaces.Task) {
	fmt.Printf("----------------------------------------------------------------------------------------------")
	fmt.Printf("Crawled url :\t" + items.GetRequest().GetUrl() + "\n")
	fmt.Printf("Crawled result : ")
	for key, value := range items.GetAll() {
		fmt.Printf(key + "\t:\t" + value)
	}
}
