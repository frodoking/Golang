package page_processor

import (
	"go_spider/core/common/page"
)

type PageProcessor interface {
	Process(p *page.Page)
}
