package request

import (
	"io/ioutil"
	"net/http"
	"os"
)

import (
	"github.com/bitly/go-simplejson"
)

import (
	"go_spider/core/common/mlog"
)

type Request struct {
	url string
	// Responce type: html json jsonp text
	respType string
	method   string
	postdata string
	// name for marking url and distinguish different urls in PageProcesser and Pipeline
	urltag    string
	header    http.Header
	cookies   []*http.Cookie
	proxyHost string
	// Redirect function for downloader used in http.Client
	// If CheckRedirect returns an error, the Client's Get
	// method returns both the previous Response.
	// If CheckRedirect returns error.New("normal"), the error process after client.Do will ignore the error.
	checkRedirect func(req *http.Request, via []*http.Request) error

	meta interface{}
}

func NewRequest(url string, respType string, urltag string, method string,
	postdata string, header http.Header, cookies []*http.Cookie,
	checkRedirect func(req *http.Request, via []*http.Request) error,
	meta interface{}) *Request {
	return &Request{url, respType, method, postdata, urltag, header, cookies, "", checkRedirect, meta}
}

func NewRequestWithProxy(url string, respType string, urltag string, method string,
	postdata string, header http.Header, cookies []*http.Cookie, proxyHost string,
	checkRedirect func(req *http.Request, via []*http.Request) error,
	meta interface{}) *Request {
	return &Request{url, respType, method, postdata, urltag, header, cookies, proxyHost, checkRedirect, meta}
}

func NewRequestWithHeaderFile(url string, respType string, headerFile string) *Request {
	_, err := os.Stat(headerFile)
	if err != nil {
		return NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
	}

	h := readHeaderFromFile(headerFile)
	return NewRequest(url, respType, "", "GET", "", h, nil, nil, nil)
}

func readHeaderFromFile(headerFile string) http.Header {
	b, err := ioutil.ReadFile(headerFile)
	if err != nil {
		mlog.LogInst().LogError(err.Error())
		return nil
	}

	js, _ := simplejson.NewJson(b)

	h := make(http.Header)
	h.Add("User-Agent", js.Get("User-Agent").MustString())
	h.Add("Referer", js.Get("Referer").MustString())
	h.Add("Cookie", js.Get("Cookie").MustString())
	h.Add("Cache-Control", "max-arg=0")
	h.Add("Connection", "keep-alive")
	return h
}

func (this *Request) AddHeaderFile(headerFile string) *Request {
	_, err := os.Stat(headerFile)
	if err != nil {
		return this
	}
	h := readHeaderFromFile(headerFile)
	this.header = h
	return this
}

func (this *Request) AddProxyHost(host string) *Request {
	this.proxyHost = host
	return this
}

func (this *Request) GetUrl() string {
	return this.url
}

func (this *Request) GetUrlTag() string {
	return this.urltag
}

func (this *Request) GetMethod() string {
	return this.method
}

func (this *Request) GetPostdata() string {
	return this.postdata
}

func (this *Request) GetHeader() http.Header {
	return this.header
}

func (this *Request) GetCookies() []*http.Cookie {
	return this.cookies
}

func (this *Request) GetProxyHost() string {
	return this.proxyHost
}

func (this *Request) GetResponseType() string {
	return this.respType
}

func (this *Request) GetRedirectFunc() func(req *http.Request, via []*http.Request) error {
	return this.checkRedirect
}

func (this *Request) GetMeta() interface{} {
	return this.meta
}
