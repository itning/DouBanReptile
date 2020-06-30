package xpath

import (
	"bytes"
	"errors"
	"github.com/antchfx/htmlquery"
	"github.com/itning/DouBanReptile/internal/error2"
	"golang.org/x/net/html"
)

type Data struct {
	Body  []byte
	Xpath string
}

type Nodes []*html.Node

func Parser(data Data) Nodes {
	if nil == data.Body || "" == data.Xpath {
		if handlerError(errors.New("data attrs must not nil\n")) {
			return nil
		}
	}
	node, e := htmlquery.Parse(bytes.NewReader(data.Body))
	handlerError(e)
	return htmlquery.Find(node, data.Xpath) // "//a//@href"
}

func (n Nodes) Text() []string {
	vs := make([]string, 0)
	for _, a := range n {
		text := htmlquery.InnerText(a)
		vs = append(vs, text)
	}
	return vs
}

func (n Nodes) String() string {
	texts := n.Text()
	r := "["
	for i, v := range texts {
		if i == len(texts)-1 {
			r += v
			break
		}
		r += v + ","
	}
	r += "]"
	return r
}

func (n Nodes) Attr(attr string) []string {
	vs := make([]string, 0)
	for _, a := range n {
		text := htmlquery.SelectAttr(a, attr)
		vs = append(vs, text)
	}
	return vs
}

func handlerError(e error) bool {
	if nil == e {
		return false
	} else {
		error2.GetImpl().Handler(e)
		return true
	}
}
