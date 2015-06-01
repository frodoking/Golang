package spider

import (
	"math/rand"
	"time"
)

import (
	"go_spider/core/common/mlog"
	"go_spider/core/common/page"
	"go_spider/core/common/page_items"
	"go_spider/core/common/request"
	"go_spider/core/common/resource_manage"
	"go_spider/core/downloader"
	"go_spider/core/page_processor"
	"go_spider/core/pipeline"
	"go_spider/core/scheduler"
)

type Sipder struct {
	taskname         string
	pPageProcessor   page_processor.PageProcessor
	pDownloader      downloader.Downloader
	pScheduler       scheduler.Scheduler
	pPipelines       []pipeline.Pipeline
	mc               resource_manage.ResourceManage
	threadnum        uint
	exitWhenComplete bool
	startSleeptime   uint
	endSleeptime     uint
	sleeptype        string
}

func NewSpider(pageinst page_processor.PageProcessor, taskname string) *Sipder {
	mlog.StraceInst().Open()
	ap := &Sipder{taskname: taskname, pPageProcessor: pageinst}

	ap.CloseFileLog()
	ap.exitWhenComplete = true
	ap.sleeptype = "fixed"
	ap.startSleeptime = 0

	if ap.pScheduler == nil {
		ap.SetScheduler(scheduler.NewQueueScheduler(false))
	}

	if ap.pDownloader == nil {
		ap.SetDownloader(downloader.NewHttpDownloader())
	}

	mlog.StraceInst().Println("** start spider **")

	ap.pPipelines = make([]pipeline.Pipeline, 0)

	return ap
}

func (this *Sipder) Taskname() string {
	return this.taskname
}

func (this *Sipder) CloseFileLog() *Sipder {
	mlog.InitFilelog(false, "")
	return this
}

func (this *Sipder) SetScheduler(s scheduler.Scheduler) *Sipder {
	this.pScheduler = s
	return this
}

func (this *Sipder) SetDownloader(d downloader.Downloader) *Sipder {
	this.pDownloader = d
	return this
}

func (this *Sipder) AddUrl(url string, respType string) *Sipder {
	req := request.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
	this.AddRequest(req)
	return this
}

func (this *Sipder) AddRequest(req *request.Request) *Sipder {
	if req == nil {
		mlog.LogInst().LogError("request is nil")
		return this
	} else if req.GetUrl() == "" {
		mlog.LogInst().LogError("request is empty")
		return this
	} else {
		this.pScheduler.Push(req)
		return this
	}
}

func (this *Sipder) AddPipeline(p pipeline.Pipeline) *Sipder {
	this.pPipelines = append(this.pPipelines, p)
	return this
}

func (this *Sipder) SetThreadnum(i uint) *Sipder {
	this.threadnum = i
	return this
}

func (this *Sipder) Run() {
	if this.threadnum == 0 {
		this.threadnum = 1
	}

	this.mc = resource_manage.NewResourceManageChan(this.threadnum)

	for {
		req := this.pScheduler.Poll()

		if this.mc.Has() == 0 && req == nil && this.exitWhenComplete {
			mlog.StraceInst().Println("** end spider **")
			break
		} else if req == nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		this.mc.GetOne()

		go func(req *request.Request) {
			defer this.mc.FreeOne()
			mlog.StraceInst().Println("start crawl : " + req.GetUrl())
			this.pageProcess(req)
		}(req)
	}

	this.close()
}

// core processer
func (this *Sipder) pageProcess(req *request.Request) {
	var p *page.Page
	defer func() {
		if err := recover(); err != nil {
			if strerr, ok := err.(string); ok {
				mlog.LogInst().LogError(strerr)
			} else {
				mlog.LogInst().LogError("pageProcess error")
			}
		}
	}()

	// download page
	for i := 0; i < 3; i++ {
		this.sleep()
		p = this.pDownloader.Download(req)
		if p.IsSucc() {
			break
		}
	}

	if !p.IsSucc() {
		return
	}

	this.pPageProcessor.Process(p)
	for _, req := range p.GetTargetRequests() {
		this.AddRequest(req)
	}

	// output
	if !p.GetSkip() {
		for _, pip := range this.pPipelines {
			pip.Process(p.GetPageItems(), this)
		}
	}

}

func (this *Sipder) sleep() {
	if this.sleeptype == "fixed" {
		time.Sleep(time.Duration(this.startSleeptime) * time.Millisecond)
	} else if this.sleeptype == "rand" {
		sleeptime := rand.Intn(int(this.endSleeptime-this.startSleeptime)) + int(this.startSleeptime)
		time.Sleep(time.Duration(sleeptime) * time.Millisecond)
	}
}

func (this *Sipder) close() {
	this.SetScheduler(scheduler.NewQueueScheduler(false))
	this.SetDownloader(downloader.NewHttpDownloader())
	this.pPipelines = make([]pipeline.Pipeline, 0)
	this.exitWhenComplete = true
}

// Deal with one url and return the PageItems with other setting.
func (this *Sipder) GetByRequest(req *request.Request) *page_items.PageItems {
	var reqs []*request.Request
	reqs = append(reqs, req)
	items := this.GetAllByRequest(reqs)
	if len(items) != 0 {
		return items[0]
	}
	return nil
}

func (this *Sipder) GetAllByRequest(reqs []*request.Request) []*page_items.PageItems {
	for _, req := range reqs {
		this.AddRequest(req)
	}

	pip := pipeline.NewCollectPipelinePageItems()
	this.AddPipeline(pip)

	this.Run()

	return pip.GetCollected()
}
