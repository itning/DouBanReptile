package scheduler

import (
	"DouBanReptile/internal/bloom"
	"DouBanReptile/internal/request"
	"DouBanReptile/internal/xpath"
	"fmt"
	"strings"
	"sync"
	"time"
)

type Dispatcher struct {
	BaseUrl string
	Headers map[string]string
	Cookies map[string]string
}

type requestData struct {
	rd       request.Data
	xpath    string
	callBack func(xpath.Nodes, request.Data)
}

type Pagination struct {
	StartPage int
	EndPage   int
}

type PaginationRange struct {
	StartSize int
	EndSize   int
	EveryAdd  int
}

var requestChan chan requestData
var wg sync.WaitGroup
var filter = bloom.NewBloomFilter()

func (d *Dispatcher) Init(firstUrl string, xpath string, callBack func(xpath.Nodes, request.Data), requestLimit time.Duration, page *Pagination) {
	requestChan = make(chan requestData)
	go d.do(requestLimit)
	if page == nil {
		d.Add(firstUrl, xpath, callBack)
	} else {
		for i := page.StartPage; i <= page.EndPage; i++ {
			d.Add(fmt.Sprintf(firstUrl, i), xpath, callBack)
		}
	}
	wg.Wait()
}

func (d *Dispatcher) Init2(firstUrl string, xpath string, callBack func(xpath.Nodes, request.Data), requestLimit time.Duration, page *PaginationRange) {
	requestChan = make(chan requestData)
	go d.do(requestLimit)
	if page == nil {
		d.Add(firstUrl, xpath, callBack)
	} else {
		for i := page.StartSize; i <= page.EndSize; i += page.EveryAdd {
			d.Add(fmt.Sprintf(firstUrl, i), xpath, callBack)
		}
	}
	wg.Wait()
}

func (d *Dispatcher) Add(url string, xpath string, callBack func(xpath.Nodes, request.Data)) {
	if strings.HasPrefix(url, d.BaseUrl) {
		url = string([]rune(url)[len(d.BaseUrl):])
	}
	if d.checkUrl(url) {
		return
	}
	wg.Add(1)
	requestChan <- requestData{rd: request.Data{Headers: d.Headers, Cookies: d.Cookies, Url: d.urlFormat(url)}, xpath: xpath, callBack: callBack}
}

func (d *Dispatcher) urlFormat(url string) string {
	// http://xxx.com/
	if strings.HasSuffix(d.BaseUrl, "/") {
		if strings.HasPrefix(url, "/") {
			return d.BaseUrl + url[1:]
		} else {
			return d.BaseUrl + url
		}
	} else {
		if strings.HasPrefix(url, "/") {
			return d.BaseUrl + url
		} else {
			return fmt.Sprintf("%s/%s", d.BaseUrl, url)
		}
	}
}

func (d *Dispatcher) do(requestLimit time.Duration) {
	for {
		time.Sleep(requestLimit)
		requestData := <-requestChan
		go d.doRequestAndParser(requestData)
	}
}

func (d *Dispatcher) doRequestAndParser(data requestData) {
	bytes := request.Handler(data.rd)
	nodes := xpath.Parser(xpath.Data{Body: bytes, Xpath: data.xpath})
	data.callBack(nodes, data.rd)
	wg.Done()
}

func (d *Dispatcher) checkUrl(url string) bool {
	if filter.Contains(url) {
		return true
	} else {
		filter.Add(url)
		if strings.HasPrefix(url, "javascript") || strings.HasPrefix(url, "/javascript") || strings.TrimSpace(url) == "" {
			return true
		}
		return false
	}
}
