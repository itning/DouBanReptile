package gui

import (
	"DouBanReptile/internal/log"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
	"os"
	"strconv"
	"strings"
)

var application fyne.App
var msgLabel *widget.Label

type Preference struct {
	GroupEntityURL             string
	MaxPrice                   int
	IncludeNoContentPriceCheck bool
	ExcludeKeyArray            []string
	MaxPage                    int
}

func (p Preference) String() string {
	return fmt.Sprintf("群组链接：%s\n最大价格：%d\n爬取不带价格的：%t\n排除关键字：%s\n爬取最大页数：%d\n",
		p.GroupEntityURL, p.MaxPrice, p.IncludeNoContentPriceCheck, p.ExcludeKeyArray, p.MaxPage)
}

func Open(onStart func(p Preference)) {
	_ = os.Setenv("FYNE_FONT", "C:\\Windows\\Fonts\\simsun.ttc")
	_ = os.Setenv("FYNE_THEME", "light")
	defer os.Unsetenv("FYNE_THEME")
	defer os.Unsetenv("FYNE_FONT")
	log.SetImpl(Log{})

	p := Preference{
		GroupEntityURL:             "/group/554566/discussion?start=%d",
		MaxPrice:                   1500,
		IncludeNoContentPriceCheck: false,
		ExcludeKeyArray:            []string{},
		MaxPage:                    10,
	}
	application = app.New()

	w := application.NewWindow("豆瓣租房小组爬虫")
	w.Resize(fyne.Size{
		Width:  400,
		Height: 200,
	})
	w.CenterOnScreen()
	groupUrlEntry := widget.NewEntry()
	groupUrlEntry.Text = p.GroupEntityURL

	maxPriceEntry := widget.NewEntry()
	maxPriceEntry.Text = strconv.Itoa(p.MaxPrice)
	maxPriceEntry.OnChanged = handlePriceInputChange(maxPriceEntry, &p)

	maxPageEntry := widget.NewEntry()
	maxPageEntry.Text = strconv.Itoa(p.MaxPage)
	maxPageEntry.OnChanged = handlePageInputChange(maxPageEntry, &p)

	excludeKeyEntry := widget.NewEntry()
	excludeKeyEntry.Text = "限女"

	isIncludeNoContentPriceCheck := widget.NewCheck("爬取不带价格的", func(b bool) {
		p.IncludeNoContentPriceCheck = b
	})

	w.SetContent(widget.NewVBox(
		widget.NewLabel("设置豆瓣群组链接："),
		groupUrlEntry,
		widget.NewLabel("设置爬取页数："),
		maxPageEntry,
		widget.NewLabel("设置最大价格："),
		maxPriceEntry,
		widget.NewLabel("设置排除关键字（用|分隔）："),
		excludeKeyEntry,
		isIncludeNoContentPriceCheck,
		widget.NewButton("开始爬取", func() {
			handleKey(excludeKeyEntry, &p)
			dialog.ShowConfirm("确认", p.String(), func(b bool) {
				if b {
					start(p, onStart)
				}
			}, w)
		}),
	))

	w.ShowAndRun()
}

func start(p Preference, onStart func(p Preference)) {
	window := application.NewWindow("爬取中...")
	window.Resize(fyne.Size{
		Width:  400,
		Height: 200,
	})
	window.CenterOnScreen()
	msgLabel = widget.NewLabel("")
	window.SetContent(widget.NewVScrollContainer(msgLabel))
	window.Show()
	onStart(p)
}

func Print(content string) {
	msgLabel.SetText(msgLabel.Text + content)
}

func handleKey(excludeKeyEntry *widget.Entry, p *Preference) {
	excludeKeyArray := strings.Split(excludeKeyEntry.Text, "|")
	p.ExcludeKeyArray = []string{}
	for _, key := range excludeKeyArray {
		if key != "" {
			p.ExcludeKeyArray = append(p.ExcludeKeyArray, key)
		}
	}
}

func handlePriceInputChange(priceEntity *widget.Entry, p *Preference) func(string) {
	return func(input string) {
		if value, err := strconv.Atoi(input); err != nil {
			if len(input) == 0 {
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
			if len(input) == 0 {
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
