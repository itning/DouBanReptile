package main

import "fmt"

func main() {
	str := "标题：滨江区 近西兴地铁口 春波南苑精装修 朝南大主卧 带内独厨独卫 独立阳台 采光好 空间大 家电齐全 有需要的联系"
	fmt.Println(string([]rune(str)[3:]))
}
