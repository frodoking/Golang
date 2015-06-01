package downloader

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html/charset"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/bitly/go-simplejson"
)

import (
	"go_spider/core/common/mlog"
	"go_spider/core/common/page"
	"go_spider/core/common/request"
	"go_spider/core/common/util"
)

// The HttpDownloader download page by package net/http.
// The "html" content is contained in dom parser of package goquery.
// The "json" content is saved.
// The "jsonp" content is modified to json.
// The "text" content will save body plain text only.
// The page result is saved in Page.
type HttpDownloader struct {
}

func NewHttpDownloader() *HttpDownloader {
	return &HttpDownloader{}
}

func (this *HttpDownloader) Download(req *request.Request) *page.Page {
	var mtype string
	var p = page.NewPage(req)
	mtype = req.GetResponseType()
	switch mtype {
	case "html":
		return this.downloadHtml(p, req)
	case "json":
		return this.downloadJson(p, req)
	case "text":
		return this.downloadText(p, req)
	default:
		mlog.LogInst().LogError("error request type:" + mtype)
	}

	return p
}

func (this *HttpDownloader) downloadHtml(p *page.Page, req *request.Request) *page.Page {
	var err error
	p, destbody := this.downloadFile(p, req)
	fmt.Printf("Destbody %v \r\n", destbody)

	if !p.IsSucc() {
		fmt.Print("Page error \r\n")
		return p
	}

	bodyReader := bytes.NewReader([]byte(destbody))
	var doc *goquery.Document
	if doc, err = goquery.NewDocumentFromReader(bodyReader); err != nil {
		mlog.LogInst().LogError(err.Error())
		p.SetStatus(true, err.Error())
		return p
	}

	var body string
	if body, err = doc.Html(); err != nil {
		mlog.LogInst().LogError(err.Error())
		p.SetStatus(true, err.Error())
		return p
	}

	p.SetBodyStr(body).SetHtmlParser(doc).SetStatus(false, "")

	return p
}

func (this *HttpDownloader) downloadJson(p *page.Page, req *request.Request) *page.Page {
	var err error
	p, destbody := this.downloadFile(p, req)
	fmt.Printf("Destbody %v \r\n", destbody)

	if !p.IsSucc() {
		fmt.Print("Page error \r\n")
		return p
	}

	var body []byte
	body = []byte(destbody)
	mtype := req.GetResponseType()
	if mtype == "jsonp" {
		tmpstr := util.JsonpToJson(destbody)
		body = []byte(tmpstr)
	}

	var r *simplejson.Json
	if r, err = simplejson.NewJson(body); err != nil {
		mlog.LogInst().LogError(string(body) + "\t" + err.Error())
		p.SetStatus(true, err.Error())
		return p
	}

	p.SetBodyStr(string(body)).SetJson(r).SetStatus(false, "")

	return p
}

func (this *HttpDownloader) downloadText(p *page.Page, req *request.Request) *page.Page {
	p, destbody := this.downloadFile(p, req)
	if !p.IsSucc() {
		return p
	}

	p.SetBodyStr(destbody).SetStatus(false, "")
	return p
}

func (this *HttpDownloader) downloadFile(p *page.Page, req *request.Request) (*page.Page, string) {
	var err error
	var urlstr string
	if urlstr = req.GetUrl(); len(urlstr) == 0 {
		mlog.LogInst().LogError("url is empty")
		p.SetStatus(true, "url is empty")
		return p, ""
	}

	var resp *http.Response

	if proxystr := req.GetProxyHost(); len(proxystr) != 0 {
		fmt.Print("HttpProxy Enter ", proxystr, "\n")
		resp, err = connectByHttpProxy(p, req)
	} else {
		fmt.Print("Http Normal Enter \n", proxystr, "\n")
		resp, err = connectByHttp(p, req)
	}

	if err != nil {
		return p, ""
	}

	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Resp body %v \r\n", string(b))

	p.SetHeader(resp.Header)
	p.SetCookies(resp.Cookies())

	bodyStr := this.changeCharsetEncodingAuto(resp.Header.Get("Content-Type"), resp.Body)
	fmt.Printf("utf-8 body %v \r\n", bodyStr)

	defer resp.Body.Close()

	return p, bodyStr
}

func connectByHttp(p *page.Page, req *request.Request) (*http.Response, error) {
	client := &http.Client{CheckRedirect: req.GetRedirectFunc()}

	httpreq, err := http.NewRequest(req.GetMethod(), req.GetUrl(), strings.NewReader(req.GetPostdata()))

	if header := req.GetHeader(); header != nil {
		httpreq.Header = req.GetHeader()
	}

	if cookies := req.GetCookies(); cookies != nil {
		for i := range cookies {
			httpreq.AddCookie(cookies[i])
		}
	}

	var resp *http.Response
	if resp, err = client.Do(httpreq); err != nil {
		if e, ok := err.(*url.Error); ok && e.Err != nil && e.Err.Error() == "normal" {
			// normal
		} else {
			mlog.LogInst().LogError(err.Error())
			p.SetStatus(true, err.Error())
			fmt.Printf("client do error %v \r\n", err)
			return nil, err
		}
	}
	return resp, nil
}

// choose a proxy server to excute http GET/method to download
func connectByHttpProxy(p *page.Page, req *request.Request) (*http.Response, error) {
	request, _ := http.NewRequest("GET", req.GetUrl(), nil)
	proxy, err := url.Parse(req.GetProxyHost())
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Charset auto determine. Use golang.org/x/net/html/charset. Get page body and change it to utf-8
func (this *HttpDownloader) changeCharsetEncodingAuto(contentTypeStr string, sor io.ReadCloser) string {
	var err error
	destReader, err := charset.NewReader(sor, contentTypeStr)

	if err != nil {
		mlog.LogInst().LogError(err.Error())
		destReader = sor
	}

	var sorbody []byte
	if sorbody, err = ioutil.ReadAll(destReader); err != nil {
		mlog.LogInst().LogError(err.Error())
	}

	bodystr := string(sorbody)
	return bodystr
}
