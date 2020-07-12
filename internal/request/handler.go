package request

import (
	"errors"
	"github.com/itning/DouBanReptile/internal/error2"
	"github.com/itning/DouBanReptile/internal/log"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	UserAgentPCChrome    = `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36`
	UserAgentPhoneChrome = `Mozilla/5.0 (Linux; Android 8.0.0; Pixel 2 XL Build/OPD1.170816.004) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Mobile Safari/537.36`
)

type Data struct {
	Headers map[string]string //optional
	Method  string            //default get
	Url     string
	Cookies map[string]string //optional
}

// use F12 dev tools in browser and write [document.cookie] in console
func AnalysisCookieString(cookies string) map[string]string {
	cookieMap := make(map[string]string)
	if "" == strings.TrimSpace(cookies) {
		return cookieMap
	}
	cookieArray := strings.Split(cookies, "; ")
	for _, cookie := range cookieArray {
		cookie := strings.Split(cookie, "=")
		if 1 == len(cookie) {
			continue
		}
		cookieMap[cookie[0]] = cookie[1]
	}
	return cookieMap
}

func (d *Data) format() error {
	if "" == d.Method {
		d.Method = http.MethodGet
	}
	if "" == d.Url {
		return errors.New("url must not be empty")
	}
	return nil
}

func (d *Data) addCookies(request *http.Request) {
	if nil == d.Cookies {
		return
	}
	for k, v := range d.Cookies {
		cookie := &http.Cookie{Name: k, Value: v}
		request.AddCookie(cookie)
	}
}

func (d *Data) addHeaders(request *http.Request) {
	if nil == d.Headers {
		return
	}
	for k, v := range d.Headers {
		request.Header.Add(k, v)
	}
}

// request handler
func Handler(data Data) []byte {
	err := data.format()
	if handlerError(err) {
		return nil
	}
	log.GetImpl().Printf("<==Method: %s Request: %s", data.Method, data.Url)
	request, e := http.NewRequest(data.Method, data.Url, nil)
	if handlerError(e) {
		return nil
	}
	data.addCookies(request)
	data.addHeaders(request)
	cli := http.Client{Timeout: time.Second * 10}
	response, e := cli.Do(request)
	if handlerError(e) {
		return nil
	}
	log.GetImpl().Printf("==>Request: %s Done With Response Status: %d", data.Url, response.StatusCode)
	readCloser := response.Body
	defer func() {
		handlerError(readCloser.Close())
	}()
	all, e := ioutil.ReadAll(readCloser)
	if handlerError(e) {
		return nil
	}
	return all
}

func handlerError(e error) bool {
	if nil == e {
		return false
	} else {
		error2.GetImpl().Handler(e)
		return true
	}
}
