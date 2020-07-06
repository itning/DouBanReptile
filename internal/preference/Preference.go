package preference

import "fmt"

type Preference struct {
	GroupEntityURL             string   // 群组URL
	MaxPrice                   int      // 最大价格
	IncludeNoContentPriceCheck bool     // 包含标题没有写价格的
	ExcludeKeyArray            []string // 排除关键字
	IncludeKeyArray            []string // 包含关键字
	MaxPage                    int      // 爬取最大页数
	SavePreference             bool     // 是否持久化配置
}

func (p Preference) String() string {
	return fmt.Sprintf("群组链接：%s\n最大价格：%d\n爬取不带价格的：%t\n爬取关键字：%s\n排除关键字：%s\n爬取最大页数：%d\n",
		p.GroupEntityURL, p.MaxPrice, p.IncludeNoContentPriceCheck, p.IncludeKeyArray, p.ExcludeKeyArray, p.MaxPage)
}
