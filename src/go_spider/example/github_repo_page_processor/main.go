package main

import (
	"fmt"
	"strings"
)

import (
	"github.com/PuerkitoBio/goquery"
)

import (
	"go_spider/core/common/page"
	"go_spider/core/pipeline"
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

	var urls []string
	query.Find("h3[class='repo-list-name'] a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		urls = append(urls, "http://github.com/"+href)
	})

	p.AddTargetRequests(urls, "html")

	name := query.Find(".entry-title .author").Text()
	name = strings.Trim(name, " \t\n")
	repository := query.Find(".entity-title .js-current-repository").Text()
	repository = strings.Trim(repository, " \t\n")

	if name == "" {
		p.SetSkip(true)
	}

	p.AddField("author", name)
	p.AddField("project", repository)
}

func main() {
	spider.NewSpider(NewMyPageProcessor(), "TaskName").
		AddUrl("https://github.com/hu17889?tab=repositories", "html").
		AddPipeline(pipeline.NewPipelineConsole()).
		SetThreadnum(3).
		Run()
}
