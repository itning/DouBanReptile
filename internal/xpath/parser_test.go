package xpath

import (
	"fmt"
	"github.com/itning/DouBanReptile/internal/request"
	"testing"
)

func TestParser(t *testing.T) {
	data := request.Data{Url: "https://www.nowcoder.com/contestRoom"}
	bytes := request.Handler(data)
	nodes := Parser(Data{Body: bytes, Xpath: `//div[@class="pagination"]//li//a`})
	for i, vv := range nodes.Attr("href") {
		fmt.Printf("%d %s\n", i+1, vv)
	}
	for i, vv := range nodes.Text() {
		fmt.Printf("%d %s\n", i+1, vv)
	}
}
