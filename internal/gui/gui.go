package gui

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
	"github.com/itning/DouBanReptile/internal/error2"
	"github.com/itning/DouBanReptile/internal/log"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var version = "1.1.1"
var author = "itning"
var application fyne.App
var msgLabel *widget.Label
var mainWindow fyne.Window
var container *widget.ScrollContainer

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

func Open(onStart func(p Preference)) {
	_ = os.Setenv("FYNE_FONT", "C:\\Windows\\Fonts\\simsun.ttc")
	_ = os.Setenv("FYNE_THEME", "light")
	defer os.Unsetenv("FYNE_THEME")
	defer os.Unsetenv("FYNE_FONT")
	error2.SetImpl(ErrorHandler{})
	log.SetImpl(Log{})

	p := Preference{
		GroupEntityURL:             "/group/554566/discussion?start=%d",
		MaxPrice:                   1500,
		IncludeNoContentPriceCheck: false,
		ExcludeKeyArray:            []string{},
		MaxPage:                    10,
	}
	application = app.New()

	mainWindow = application.NewWindow("豆瓣租房小组爬虫 ver:" + version + " by " + author)
	mainWindow.Resize(fyne.Size{
		Width:  400,
		Height: 200,
	})
	mainWindow.CenterOnScreen()

	githubURL, _ := url.Parse("https://github.com/itning/DouBanReptile")
	hyperLink := widget.NewHyperlink("访问GitHub点个Star", githubURL)

	groupUrlEntry := widget.NewEntry()
	groupUrlEntry.Text = p.GroupEntityURL
	groupUrlEntry.OnChanged = func(s string) {
		groupUrl := strings.TrimSpace(s)
		if "" == groupUrl {
			return
		}
		p.GroupEntityURL = groupUrl
	}

	maxPriceEntry := widget.NewEntry()
	maxPriceEntry.Text = strconv.Itoa(p.MaxPrice)
	maxPriceEntry.OnChanged = handlePriceInputChange(maxPriceEntry, &p)

	maxPageEntry := widget.NewEntry()
	maxPageEntry.Text = strconv.Itoa(p.MaxPage)
	maxPageEntry.OnChanged = handlePageInputChange(maxPageEntry, &p)

	includeKeyEntry := widget.NewEntry()
	includeKeyEntry.Text = ""

	excludeKeyEntry := widget.NewEntry()
	excludeKeyEntry.Text = "限女"

	isIncludeNoContentPriceCheck := widget.NewCheck("也爬取不带价格的", func(b bool) {
		p.IncludeNoContentPriceCheck = b
	})

	isSavePreferenceCheck := widget.NewCheck("保存配置(写入EXE所在目录DouBanConfig.ini文件)", func(b bool) {
		p.SavePreference = b
	})
	isSavePreferenceCheck.Checked = true

	mainWindow.SetContent(widget.NewVBox(
		hyperLink,
		widget.NewLabel("设置豆瓣群组链接："),
		groupUrlEntry,
		widget.NewLabel("设置爬取页数："),
		maxPageEntry,
		widget.NewLabel("设置最大价格："),
		maxPriceEntry,
		widget.NewLabel("设置标题爬取关键字（用|分隔，不写全爬）："),
		includeKeyEntry,
		widget.NewLabel("设置标题和内容排除关键字（用|分隔）："),
		excludeKeyEntry,
		isIncludeNoContentPriceCheck,
		isSavePreferenceCheck,
		widget.NewButton("开始爬取", func() {
			p.IncludeKeyArray = handleKey(includeKeyEntry)
			p.ExcludeKeyArray = handleKey(excludeKeyEntry)
			dialog.ShowConfirm("确认", p.String(), func(b bool) {
				if b {
					start(p, onStart)
				}
			}, mainWindow)
		}),
	))

	mainWindow.ShowAndRun()
}

func closeMainWindow() {
	mainWindow.Close()
}

func start(p Preference, onStart func(p Preference)) {
	window := application.NewWindow("爬取中...")
	window.Resize(fyne.Size{
		Width:  500,
		Height: 200,
	})
	window.CenterOnScreen()
	msgLabel = widget.NewLabel("")
	container = widget.NewVScrollContainer(msgLabel)
	window.SetContent(container)
	window.Show()
	closeMainWindow()
	window.SetOnClosed(func() {
		application.Quit()
	})
	go onStart(p)
}

func Print(content string) {
	msgLabel.SetText(msgLabel.Text + content)
	adjust := msgLabel.MinSize().Height - container.Size().Height
	container.Offset = fyne.NewPos(0, adjust)
}

func handleKey(excludeKeyEntry *widget.Entry) []string {
	excludeKeyArray := strings.Split(excludeKeyEntry.Text, "|")
	var temp []string
	for _, key := range excludeKeyArray {
		if key != "" {
			temp = append(temp, key)
		}
	}
	return temp
}

func handlePriceInputChange(priceEntity *widget.Entry, p *Preference) func(string) {
	return func(input string) {
		if value, err := strconv.Atoi(input); err != nil {
			if 0 == len(input) {
				priceEntity.SetText("")
				p.MaxPrice = 9999
			} else {
				priceEntity.SetText(input[:len(input)-1])

			}
		} else {
			p.MaxPrice = value
		}
	}
}

func handlePageInputChange(pageEntity *widget.Entry, p *Preference) func(string) {
	return func(input string) {
		if value, err := strconv.Atoi(input); err != nil {
			if 0 == len(input) {
				pageEntity.SetText("1")
				p.MaxPage = 1
			} else {
				pageEntity.SetText(input[:len(input)-1])
			}
		} else {
			if value < 1 {
				pageEntity.SetText("1")
				p.MaxPage = 1
			} else {
				p.MaxPage = value
			}
		}
	}
}
