package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
	"github.com/itning/DouBanReptile/internal/error2"
	"github.com/itning/DouBanReptile/internal/ini"
	"github.com/itning/DouBanReptile/internal/log"
	"github.com/itning/DouBanReptile/internal/preference"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var version = "1.1.5"
var author = "itning"
var application fyne.App
var msgLabel *widget.Label
var mainWindow fyne.Window
var container *widget.ScrollContainer

func Open(onStart func(p preference.Preference)) {
	_ = os.Setenv("FYNE_FONT", "C:\\Windows\\Fonts\\simsun.ttc")
	_ = os.Setenv("FYNE_THEME", "light")
	defer os.Unsetenv("FYNE_THEME")
	defer os.Unsetenv("FYNE_FONT")
	error2.SetImpl(ErrorHandler{})
	log.SetImpl(Log{})
	config := ini.Config{}
	p := config.Read()
	checkPreference(p)

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

	cookieEntry := widget.NewEntry()
	cookieEntry.MultiLine = true
	cookieEntry.Text = splitCookieStringOnGUI(p.CookieString)
	cookieEntry.OnChanged = func(s string) {
		cookieEntry.Text = splitCookieStringOnGUI(s)
		p.CookieString = s
	}

	maxPriceEntry := widget.NewEntry()
	maxPriceEntry.Text = strconv.Itoa(p.MaxPrice)
	maxPriceEntry.OnChanged = handlePriceInputChange(maxPriceEntry, p)

	maxPageEntry := widget.NewEntry()
	maxPageEntry.Text = strconv.Itoa(p.MaxPage)
	maxPageEntry.OnChanged = handlePageInputChange(maxPageEntry, p)

	includeKeyEntry := widget.NewEntry()
	includeKeyEntry.Text = strings.Join(p.IncludeKeyArray, "|")

	excludeKeyEntry := widget.NewEntry()
	excludeKeyEntry.Text = strings.Join(p.ExcludeKeyArray, "|")

	isIncludeNoContentPriceCheck := widget.NewCheck("也爬取不带价格的", func(b bool) {
		p.IncludeNoContentPriceCheck = b
	})

	isSavePreferenceCheck := widget.NewCheck("保存配置(写入EXE所在目录DouBanConfig.ini文件)", func(b bool) {
		p.SavePreference = b
	})
	p.SavePreference = true
	isSavePreferenceCheck.Checked = true

	mainWindow.SetContent(widget.NewVBox(
		hyperLink,
		widget.NewLabel("设置豆瓣群组链接："),
		groupUrlEntry,
		widget.NewLabel("设置Cookie（如果遇到\"nil\"请设置该项）："),
		cookieEntry,
		widget.NewLabel("设置爬取页数："),
		maxPageEntry,
		widget.NewLabel("设置最大价格："),
		maxPriceEntry,
		widget.NewLabel("设置标题或内容爬取关键字（用|分隔，不写全爬）："),
		includeKeyEntry,
		widget.NewLabel("设置标题或内容排除关键字（用|分隔，不写全爬）："),
		excludeKeyEntry,
		isIncludeNoContentPriceCheck,
		isSavePreferenceCheck,
		widget.NewButton("开始爬取", func() {
			p.IncludeKeyArray = handleKey(includeKeyEntry)
			p.ExcludeKeyArray = handleKey(excludeKeyEntry)
			dialog.ShowConfirm("确认", p.String(), func(b bool) {
				if b {
					start(*p, onStart)
				}
			}, mainWindow)
		}),
	))

	mainWindow.ShowAndRun()
}

func splitCookieStringOnGUI(s string) string {
	cookie := ""
	array := strings.Split(s, ";")
	for index, item := range array {
		item = strings.TrimSpace(item)
		if index == len(array)-1 {
			cookie += item
		} else {
			cookie += item + ";\n"
		}
	}
	return cookie
}

func checkPreference(p *preference.Preference) {
	if "" == p.GroupEntityURL {
		p.GroupEntityURL = "/group/554566/discussion?start=%d"
	}
	if 0 == p.MaxPage {
		p.MaxPage = 1
	}
	if 0 == p.MaxPrice {
		p.MaxPrice = 1000
	}
}

func closeMainWindow() {
	mainWindow.Close()
}

func start(p preference.Preference, onStart func(p preference.Preference)) {
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

func handlePriceInputChange(priceEntity *widget.Entry, p *preference.Preference) func(string) {
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

func handlePageInputChange(pageEntity *widget.Entry, p *preference.Preference) func(string) {
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
