package main

import (
	"DouBanReptile/internal/markdown"
	"DouBanReptile/internal/request"
	"DouBanReptile/internal/scheduler"
	"DouBanReptile/internal/xpath"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var dispatcher scheduler.Dispatcher
var file *os.File
var priceCompile = regexp.MustCompile(`\b\d{4}\b`)
var dataArray = make(markdown.DataArray, 0)

func main() {
	headerMap := make(map[string]string)

	headerMap["User-Agent"] = request.UserAgentPCChrome

	dispatcher = scheduler.Dispatcher{
		BaseUrl: "https://www.douban.com",
		Headers: headerMap,
	}
	dispatcher.Init2(
		"/group/554566/discussion?start=%d",
		`//td[@class='title']/a`,
		each,
		time.Millisecond*500,
		&scheduler.PaginationRange{StartSize: 0, EndSize: 50, EveryAdd: 25})

	write2File()
}

func write2File() {
	var err error
	file, err = os.Create("output.md")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	write([]byte(dataArray.String()))
}

func write(b []byte) {
	if _, err := file.Write(b); err != nil {
		panic(err)
	}
}

func each(nodes xpath.Nodes, request request.Data) {
	hrefs := nodes.Attr("href")
	titles := nodes.Attr("title")
	for index, href := range hrefs {
		// 不出现"限女"
		if !strings.Contains(titles[index], "限女") {
			price := getPriceFromString(titles[index])
			// 价格1500以内
			if price != 0 && price <= 1500 {
				dispatcher.Add(href, `//div[@class="article"]`, content)
			} else if price == 0 {
				//dispatcher.Add(href, `//div[@class="article"]`, content)
			}
		}
	}
}

func getPriceFromString(title string) int {
	priceArray := priceCompile.FindAllString(title, -1)
	if len(priceArray) != 0 {
		price, e := strconv.Atoi(priceArray[0])
		if e != nil {
			log.Printf("Transform Error %s", e.Error())
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
		content := handleContent(node)
		if strings.Contains(content, "限女") {
			continue
		}
		// 处理图片
		imgArray := handleImages(node)
		// 处理时间
		timeStr := handleTime(node)

		dataArray = append(dataArray, markdown.Data{
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

func handleContent(node *html.Node) string {
	contentNode := htmlquery.FindOne(node, `//td[@class="topic-content"]`)
	if contentNode == nil {
		contentNode = htmlquery.FindOne(node, `//div[@class="topic-richtext"]`)
	}
	content := htmlquery.InnerText(contentNode)
	return content
}

func handleTitle(node *html.Node) string {
	titleNode := htmlquery.FindOne(node, `//td[@class="tablecc"]`)
	if titleNode == nil {
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