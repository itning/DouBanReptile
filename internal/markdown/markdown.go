package markdown

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"sync/atomic"
	"time"
)

var loc, _ = time.LoadLocation("Asia/Shanghai")
var num int32

type Data struct {
	TimeString string
	Time       time.Time
	Price      int
	Link       string
	Title      string
	Content    string
	Images     []string
}

type DataArray []Data

func (p *DataArray) Append(d Data) {
	*p = append(*p, d)
}

func (p DataArray) String() string {
	sort.Sort(p)
	str := ""
	for _, data := range p {
		atomic.AddInt32(&num, 1)
		str += fmt.Sprintf("%d. %s %s\n", num, data.TimeString, data.handleTitleToString())
		str += fmt.Sprintf("%s\n", data.handleContentToString())
		for _, img := range data.handleImageToString() {
			str += fmt.Sprintf("%s\n", img)
		}
	}
	return str
}

func (p DataArray) Len() int {
	return len(p)
}

func (p DataArray) Less(i, j int) bool {
	if p[i].Price != 0 && p[j].Price != 0 {
		return p[i].Price < p[j].Price
	} else {
		return p[i].Time.After(p[j].Time)
	}
}
func (p DataArray) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func String2Time(timeString string) time.Time {
	theTime, err := time.ParseInLocation("2006-01-02 15:04:05", timeString, loc)
	handlerError(err)
	return theTime
}

func (d Data) handleImageToString() []string {
	vs := make([]string, 0)
	for _, img := range d.Images {
		vs = append(vs, fmt.Sprintf(`   ![%s](%s)`, img, img))
	}
	return vs
}

func (d Data) handleContentToString() string {
	return "   " + strings.TrimSpace(d.Content)
}

func (d Data) handleTitleToString() string {
	return fmt.Sprintf("[%s](%s)", d.Title, d.Link)
}

func handlerError(e error) {
	if e != nil {
		log.Printf("Have Error %s", e.Error())
		panic(e)
	}
}
