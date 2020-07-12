package main

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/itning/DouBanReptile/internal/gui"
	"github.com/itning/DouBanReptile/internal/ini"
	"github.com/itning/DouBanReptile/internal/log"
	"github.com/itning/DouBanReptile/internal/markdown"
	"github.com/itning/DouBanReptile/internal/preference"
	"github.com/itning/DouBanReptile/internal/request"
	"github.com/itning/DouBanReptile/internal/scheduler"
	"github.com/itning/DouBanReptile/internal/xpath"
	"golang.org/x/net/html"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var dispatcher scheduler.Dispatcher
var file *os.File
var priceCompile = regexp.MustCompile(`\b\d{4}\b`) // 匹配标题中价格正则
var dataArray = make(markdown.DataArray, 0)
var pre preference.Preference

func main() {
	gui.Open(func(p preference.Preference) {
		pre = p
		savePreference(&p)
		headerMap := make(map[string]string)
		headerMap["User-Agent"] = request.UserAgentPCChrome
		headerMap["Host"] = "www.douban.com"

		dispatcher = scheduler.Dispatcher{
			BaseUrl: "https://www.douban.com",
			Headers: headerMap,
			Cookies: request.AnalysisCookieString(p.CookieString),
		}
		dispatcher.Init2(
			pre.GroupEntityURL,
			`//td[@class='title']/a`,
			each,
			time.Millisecond*500,
			&scheduler.PaginationRange{StartSize: 0, EndSize: p.MaxPage * 25, EveryAdd: 25})

		write2File()
	})
}

func savePreference(preference *preference.Preference) {
	if preference.SavePreference {
		config := ini.Config{}
		config.Write(preference)
	}
}

func write2File() {
	var err error
	file, err = os.Create("爬取结果.md")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	logger := log.GetImpl()
	logger.Printf("共爬取条数: %d", len(dataArray))
	logger.Printf("写入文件中....")
	write([]byte(dataArray.String()))
	logger.Printf("文件写入成功，请在本EXE目录下找【%s】文件", file.Name())
}

func write(b []byte) {
	if _, err := file.Write(b); err != nil {
		panic(err)
	}
}

func each(nodes xpath.Nodes, request request.Data) {
	if nil == nodes {
		log.GetImpl().Printf("Node Is Nil. Url: %s", request.Url)
		return
	}
	hrefs := nodes.Attr("href")
	titles := nodes.Attr("title")
	for index, href := range hrefs {
		if isExcludeContent(titles[index]) {
			continue
		}
		price := getPriceFromString(titles[index])
		if pre.IncludeNoContentPriceCheck {
			dispatcher.Add(href, `//div[@class="article"]`, content)
		} else if 0 != price && price <= (pre.MaxPrice) {
			// 标题上有价格并且价格在用户设置范围内
			dispatcher.Add(href, `//div[@class="article"]`, content)
		} else if 0 == price {
			// 标题上没有价格但是内容中可能有价格
			dispatcher.Add(href, `//div[@class="article"]`, contentWithTitleNoPrice)
		}
	}
}

// 只要有一个关键字存在即返回真
func isIncludeContent(content string) bool {
	// 未设置关键字，则返回真
	if 0 == len(pre.IncludeKeyArray) {
		return true
	}
	for _, key := range pre.IncludeKeyArray {
		if strings.Contains(content, key) {
			return true
		}
	}
	return false
}

// 只要关键字在标题或内容中即返回真
func checkTitleOrContentHaveKey(title string, content string) bool {
	if isIncludeContent(title) {
		return true
	} else if isIncludeContent(content) {
		return true
	} else {
		return false
	}
}

func isExcludeContent(content string) bool {
	for _, key := range pre.ExcludeKeyArray {
		if strings.Contains(content, key) {
			return true
		}
	}
	return false
}

func getPriceFromString(title string) int {
	priceArray := priceCompile.FindAllString(title, -1)
	if len(priceArray) != 0 {
		price, e := strconv.Atoi(priceArray[0])
		if e != nil {
			log.GetImpl().Printf("Transform Error %s", e.Error())
			panic(e)
		}
		return price
	} else {
		return 0
	}
}

func content(nodes xpath.Nodes, request request.Data) {
	for _, node := range nodes {
		// 处理标题
		title := handleTitle(node)
		// 处理内容
		content := handleContent(node, request)
		// 排除关键字
		if isExcludeContent(content) {
			continue
		}
		if !checkTitleOrContentHaveKey(title, content) {
			continue
		}
		// 处理图片
		imgArray := handleImages(node)
		// 处理时间
		timeStr := handleTime(node)

		dataArray.Append(markdown.Data{
			TimeString: timeStr,
			Time:       markdown.String2Time(timeStr),
			Title:      format(title),
			Price:      getPriceFromString(title),
			Link:       request.Url,
			Content:    content,
			Images:     imgArray,
		})
	}
}

func contentWithTitleNoPrice(nodes xpath.Nodes, request request.Data) {
	for _, node := range nodes {
		// 处理内容
		content := handleContent(node, request)
		price := getPriceFromString(content)
		if 0 == price || price > (pre.MaxPrice) {
			// 内容中价格依然没有或者价格大于用户设定价格
			continue
		}
		// 排除关键字
		if isExcludeContent(content) {
			continue
		}
		// 处理标题
		title := handleTitle(node)

		if !checkTitleOrContentHaveKey(title, content) {
			continue
		}
		// 处理图片
		imgArray := handleImages(node)
		// 处理时间
		timeStr := handleTime(node)

		dataArray.Append(markdown.Data{
			TimeString: timeStr,
			Time:       markdown.String2Time(timeStr),
			Title:      format(title),
			Price:      getPriceFromString(title),
			Link:       request.Url,
			Content:    content,
			Images:     imgArray,
		})
	}
}

func handleTime(node *html.Node) string {
	timeNode := htmlquery.FindOne(node, `//span[@class="color-green"]`)
	timeStr := htmlquery.InnerText(timeNode)
	return timeStr
}

func handleImages(node *html.Node) []string {
	imgs := htmlquery.Find(node, `//div[@class="image-wrapper"]//img`)
	imgArray := make([]string, 0)
	if len(imgs) != 0 {
		for _, img := range imgs {
			imgArray = append(imgArray, htmlquery.SelectAttr(img, "src"))
		}
	}
	return imgArray
}

func handleContent(node *html.Node, request request.Data) string {
	contentNode := htmlquery.FindOne(node, `//td[@class="topic-content"]`)
	if nil == contentNode {
		contentNode = htmlquery.FindOne(node, `//div[@class="topic-richtext"]`)
	}
	if nil == contentNode {
		contentNode = htmlquery.FindOne(node, `//div[@class="topic-content"]`)
	}
	if nil == contentNode {
		contentNode = htmlquery.FindOne(node, `//div[@class='rich-content topic-richtext']`)
	}
	if nil == contentNode {
		log.GetImpl().Printf("Content Is Nil So Jump Over. URL=%s", request.Url)
		return fmt.Sprintf("<内容爬取失败 URL=%s>", request.Url)
	}
	content := htmlquery.InnerText(contentNode)
	return content
}

func handleTitle(node *html.Node) string {
	titleNode := htmlquery.FindOne(node, `//td[@class="tablecc"]`)
	if nil == titleNode {
		titleNode = htmlquery.FindOne(node, `//h1`)
	}
	title := htmlquery.InnerText(titleNode)
	if strings.Contains(title, "标题") {
		title = string([]rune(title)[3:])
	}
	return title
}

func format(str string) string {
	return strings.TrimSpace(str)
}
