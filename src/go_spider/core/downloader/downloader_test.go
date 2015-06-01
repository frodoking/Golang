package downloader

import (
	"fmt"
	"testing"
)

import (
	"github.com/PuerkitoBio/goquery"
)

import (
	"go_spider/core/common/page"
	"go_spider/core/common/request"
)

func TestDownloadHtml(t *testing.T) {
	var req *request.Request
	req = request.NewRequest("http://live.sina.com.cn/zt/l/v/finance/globalnews1/", "html", "", "GET", "", nil, nil, nil, nil)
	var dl Downloader
	dl = NewHttpDownloader()

	var p *page.Page
	p = dl.Download(req)

	var doc *goquery.Document
	doc = p.GetHtmlParser()
	fmt.Println(doc)
	body := p.GetBodyStr()
	fmt.Println(body)

	var s *goquery.Selection
	s = doc.Find("body")
	if s.Length() < 1 {
		t.Error("html parse failed!")
	}
}

func TestDownloadJson(t *testing.T) {
	var req *request.Request
	req = request.NewRequest("http://live.sina.com.cn/zt/api/l/get/finance/globalnews1/index.htm?format=json&id=23521&pagesize=4&dire=f&dpc=1", "json", "", "GET", "", nil, nil, nil, nil)

	var dl Downloader
	dl = NewHttpDownloader()
	var p *page.Page
	p = dl.Download(req)

	var jsonMap interface{}
	jsonMap = p.GetJson()
	fmt.Printf("%v", jsonMap)
}

func TestCharSetChange(t *testing.T) {
	var req *request.Request
	req = request.NewRequest("http://soft.chinabyte.com/416/13164916.shtml", "html", "", "GET", "", nil, nil, nil, nil)

	var dl Downloader
	dl = NewHttpDownloader()

	var p *page.Page
	p = dl.Download(req)
	body := p.GetBodyStr()
	fmt.Println(body)
}
