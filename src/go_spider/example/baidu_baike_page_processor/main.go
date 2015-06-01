package main

import (
	"fmt"
	"strings"
)

import (
	"go_spider/core/common/page"
	"go_spider/core/common/request"
	"go_spider/core/spider"
)

type MyPageProcessor struct {
}

func NewMyPageProcessor() *MyPageProcessor {
	return &MyPageProcessor{}
}

func (this *MyPageProcessor) Process(p *page.Page) {
	if !p.IsSucc() {
		fmt.Println(p.Errormsg())
		return
	}

	query := p.GetHtmlParser()

	name := query.Find(".lemmaTitleH1").Text()
	name = strings.Trim(name, " \t\n")

	summary := query.Find(".card-summary-content .para").Text()
	summary = strings.Trim(summary, " \t\n")

	p.AddField("name", name)
	p.AddField("summary", summary)
}

func main() {
	spider := spider.NewSpider(NewMyPageProcessor(), "TaskName")
	req := request.NewRequest("http://baike.baidu.com/view/1628025.htm?fromtitle=http&fromid=243074&type=syn", "html", "", "GET", "", nil, nil, nil, nil)
	pageItems := spider.GetByRequest(req)

	url := pageItems.GetRequest().GetUrl()
	fmt.Println("----------------------- spider Get -----------------------------")
	fmt.Println("url \t: \t" + url)
	for name, value := range pageItems.GetAll() {
		fmt.Println(name + "\t:\t" + value)
	}

	fmt.Println("\n----------------------- spider GetAll -----------------------------")
	urls := []string{
		"http://baike.baidu.com/view/1628025.htm?fromtitle=http&fromid=243074&type=syn",
		"http://baike.baidu.com/view/383720.htm?fromtitle=html&fromid=97049&type=syn",
	}

	var reqs []*request.Request
	for _, url := range urls {
		req := request.NewRequest(url, "html", "", "GET", "", nil, nil, nil, nil)
		reqs = append(reqs, req)
	}

	pageItemsArr := spider.SetThreadnum(2).GetAllByRequest(reqs)
	for _, item := range pageItemsArr {
		url = item.GetRequest().GetUrl()
		fmt.Println("url\t:\t" + url)
		fmt.Printf("item\t:\t%s\n", item.GetAll())
	}
}
