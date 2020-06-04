package markdown

import (
	"fmt"
	"testing"
	"time"
)

func TestSort(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	theTime1, _ := time.ParseInLocation("2006-01-02 15:04:05", "2020-06-03 23:54:35", loc)
	theTime2, _ := time.ParseInLocation("2006-01-02 15:04:05", "2020-06-03 20:54:35", loc)
	theTime3, _ := time.ParseInLocation("2006-01-02 15:04:05", "2020-06-03 22:56:35", loc)
	theTime4, _ := time.ParseInLocation("2006-01-02 15:04:05", "2020-06-03 22:55:35", loc)
	persons := DataArray{
		{Price: 1200, Time: theTime1, Content: "哈哈", Link: "https://www.douban.com/group/topic/178647758/", Title: "as"},
		{Price: 3200, Time: theTime2, Content: "哈哈", Link: "https://www.douban.com/group/topic/178647758/", Title: "as"},
		{Price: 1600, Time: theTime3, Content: "哈哈", Link: "https://www.douban.com/group/topic/178647758/", Title: "as"},
		{Price: 1600, Time: theTime4, Content: "哈哈", Link: "https://www.douban.com/group/topic/178647758/", Title: "as"},
	}
	//persons = append(persons, Data{price: 3200, time: theTime1})
	fmt.Println(persons)
}
