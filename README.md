# 豆瓣租房爬虫

[![GitHub stars](https://img.shields.io/github/stars/itning/DouBanReptile.svg?style=social&label=Stars)](https://github.com/itning/DouBanReptile/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/itning/DouBanReptile.svg?style=social&label=Fork)](https://github.com/itning/DouBanReptile/network/members)
[![GitHub watchers](https://img.shields.io/github/watchers/itning/DouBanReptile.svg?style=social&label=Watch)](https://github.com/itning/DouBanReptile/watchers)
[![GitHub followers](https://img.shields.io/github/followers/itning.svg?style=social&label=Follow)](https://github.com/itning?tab=followers)

[![GitHub issues](https://img.shields.io/github/issues/itning/DouBanReptile.svg)](https://github.com/itning/DouBanReptile/issues)
[![GitHub license](https://img.shields.io/github/license/itning/DouBanReptile.svg)](https://github.com/itning/DouBanReptile/blob/master/LICENSE)
[![GitHub last commit](https://img.shields.io/github/last-commit/itning/DouBanReptile.svg)](https://github.com/itning/DouBanReptile/commits)
[![GitHub release](https://img.shields.io/github/release/itning/DouBanReptile.svg)](https://github.com/itning/DouBanReptile/releases)
[![GitHub repo size in bytes](https://img.shields.io/github/repo-size/itning/DouBanReptile.svg)](https://github.com/itning/DouBanReptile)
[![HitCount](http://hits.dwyl.io/itning/DouBanReptile.svg)](http://hits.dwyl.io/itning/DouBanReptile)
[![language](https://img.shields.io/badge/language-GO-green.svg)](https://github.com/itning/DouBanReptile)

## 下载

[https://github.com/itning/DouBanReptile/releases](https://github.com/itning/DouBanReptile/releases)

## 构建

```shell
go build -ldflags="-s -w -H windowsgui" -o ..\bin\main.exe DouBanReptile/cmd
```

**爬取结果文件（markdown）建议使用[typora](https://typora.io/)打开**

### 截图

![a1](https://raw.githubusercontent.com/itning/DouBanReptile/master/pic/a1.png)

![a2](https://raw.githubusercontent.com/itning/DouBanReptile/master/pic/a2.png)

![a3](https://raw.githubusercontent.com/itning/DouBanReptile/master/pic/a3png)

![a4](https://raw.githubusercontent.com/itning/DouBanReptile/master/pic/a4.png)

## 使用教程

**确保`C:\\Windows\\Fonts\\`目录下有`simsun.ttc`字体文件**

![e](https://raw.githubusercontent.com/itning/DouBanReptile/master/pic/e.png)

1. 如何设置豆瓣群组链接？

   1. 首先搜索某个地区租房，例如：`北京租房`

      ![f](https://raw.githubusercontent.com/itning/DouBanReptile/master/pic/f.png)

   2. 点进去要爬取的某个小组，例如第一个：`北京租房`

   3. 将页面拉到最下面有个`> 更多小组讨论`超链接，点进去

      ![g](https://raw.githubusercontent.com/itning/DouBanReptile/master/pic/g.png)

   4. 复制地址栏中地址（从/group开始复制到结尾），粘贴到软件`设置豆瓣群组链接`

      **有时候粘贴进软件会崩溃，不知道什么原因，建议把软件中原来的链接删除再粘贴进去。**

      ![h](https://raw.githubusercontent.com/itning/DouBanReptile/master/pic/h.png)

      ![i](https://raw.githubusercontent.com/itning/DouBanReptile/master/pic/i.png)

   5. 将`start=`后边的数字`50`改成`%d`

      ![j](https://raw.githubusercontent.com/itning/DouBanReptile/master/pic/j.png)

   6. 完成

2. 如何设置排除关键字？

   排除关键字是标题和内容只要出现关键字就会排除掉该条租房信息。

   例如默认是`限女`这个关键字，只要租房信息中包含`限女生入住`，`只限女生`等出现`限女`关键字的一律不爬。

   多个关键字用`|`分隔，注意是英文的。

   例如：`限女|短租|整租`，这三个关键字设置后，只要标题和内容出现这三个关键字软件就不会爬取。

3. 关于识别标题中的价格

   使用正则`\b\d{4}\b`识别标题中的价格信息，无法爬取少于1000元的信息。

4. 关于爬取结果排序

   先根据价格从小到大排序，价格相同根据发帖时间排序。

5. 关于爬取结果文件(.md扩展名)如何打开

   建建议下载软件：[typora](https://typora.io/)

## 测试

| 操作系统        | 测试结果 |
| --------------- | -------- |
| windows 7 sp1   | OK       |
| windows 10 1909 | OK       |