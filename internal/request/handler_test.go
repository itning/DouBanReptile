package request

import (
	"testing"
)

func TestHandler(t *testing.T) {
	data := Data{Url: "https://www.baidu.com"}
	bytes := Handler(data)
	t.Log(string(bytes))
}

func TestAnalysisCookieString(t *testing.T) {
	cookieMap := AnalysisCookieString(`NOWCODERUID=792B9A5DA08as4A826ADS2FBFACFDF9BFD55A3; NOWCODERCLINETID=674ECEE734D6C1D29455B91; gr_user_id=5d4-ffc3-44b2-84ae-0e272a`)
	t.Log(cookieMap)
}
