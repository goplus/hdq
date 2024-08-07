/*
Copyright 2021 The GoPlus Authors (goplus.org)
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package hdq

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// -----------------------------------------------------------------------------

// Printf prints the NodeSet context and `print(format, params...)`.
func (p NodeSet) Printf(w io.Writer, format string, params ...interface{}) NodeSet {
	if p.Err != nil {
		return p
	}
	p.Data.ForEach(func(node *html.Node) error {
		html.Render(w, node)
		fmt.Fprintf(w, format, params...)
		return nil
	})
	return p
}

// Dump prints the NodeSet context and `print("\n\n")`.
func (p NodeSet) Dump() NodeSet {
	return p.Printf(os.Stdout, "\n\n")
}

// -----------------------------------------------------------------------------

// ChildEqualText returns NodeSet which child node text equals `text`.
func (p NodeSet) ChildEqualText(text string) (ret NodeSet) {
	return p.Match(func(node *html.Node) bool {
		return childEqualText(node, text)
	})
}

// TextEqual returns NodeSet which node type is TextNode and it's text equals `text`.
func (p NodeSet) TextEqual(text string) (ret NodeSet) {
	return p.Match(func(node *html.Node) bool {
		return equalText(node, text)
	})
}

// TextContains returns NodeSet which node type is TextNode and it's text contains `text`.
func (p NodeSet) TextContains(text string) (ret NodeSet) {
	return p.Match(func(node *html.Node) bool {
		return containsText(node, text)
	})
}

// TextHasPrefix returns NodeSet which node type is TextNode and its prefix is `text`.
func (p NodeSet) TextHasPrefix(text string) (ret NodeSet) {
	return p.Match(func(node *html.Node) bool {
		return hasPrefixText(node, text)
	})
}

func (p NodeSet) dataAtom(elem atom.Atom) (ret NodeSet) {
	return p.Match(func(node *html.Node) bool {
		return node.DataAtom == elem
	})
}

// Element returns NodeSet which node type is ElementNode and it's element type is `elemType`.
func (p NodeSet) Element__0(elemType atom.Atom) (ret NodeSet) {
	return p.dataAtom(elemType)
}

// Element returns NodeSet which node type is ElementNode and it's element type is `elemType`.
func (p NodeSet) Element__1(elemType string) (ret NodeSet) {
	return p.Match(func(node *html.Node) bool {
		return node.Type == html.ElementNode && node.Data == elemType
	})
}

// Attribute returns NodeSet which the value of attribute `k` is `v`.
func (p NodeSet) Attribute__0(k, v string) (ret NodeSet) {
	return p.Match(func(node *html.Node) bool {
		if node.Type != html.ElementNode {
			return false
		}
		for _, attr := range node.Attr {
			if attr.Key == k && attr.Val == v {
				return true
			}
		}
		return false
	})
}

func (p NodeSet) Attribute__1(k string, filter func(v string) bool) (ret NodeSet) {
	return p.Match(func(node *html.Node) bool {
		if node.Type != html.ElementNode {
			return false
		}
		for _, attr := range node.Attr {
			if attr.Key == k && filter(attr.Val) {
				return true
			}
		}
		return false
	})
}

// ContainsClass returns NodeSet which class contains `v`.
func (p NodeSet) ContainsClass(v string) (ret NodeSet) {
	return p.Match(func(node *html.Node) bool {
		if node.Type != html.ElementNode {
			return false
		}
		for _, attr := range node.Attr {
			if attr.Key == "class" {
				return containsClass(attr.Val, v)
			}
		}
		return false
	})
}

// H1 returns NodeSet which node type is ElementNode and it's element type is `h1`.
func (p NodeSet) H1() (ret NodeSet) {
	return p.dataAtom(atom.H1)
}

// H2 returns NodeSet which node type is ElementNode and it's element type is `h2`.
func (p NodeSet) H2() (ret NodeSet) {
	return p.dataAtom(atom.H2)
}

// H3 returns NodeSet which node type is ElementNode and it's element type is `h3`.
func (p NodeSet) H3() (ret NodeSet) {
	return p.dataAtom(atom.H3)
}

// H4 returns NodeSet which node type is ElementNode and it's element type is `h4`.
func (p NodeSet) H4() (ret NodeSet) {
	return p.dataAtom(atom.H4)
}

// Td returns NodeSet which node type is ElementNode and it's element type is `td`.
func (p NodeSet) Td() (ret NodeSet) {
	return p.dataAtom(atom.Td)
}

// A returns NodeSet which node type is ElementNode and it's element type is `a`.
func (p NodeSet) A() (ret NodeSet) {
	return p.dataAtom(atom.A)
}

// P returns NodeSet which node type is ElementNode and it's element type is `p`.
func (p NodeSet) P() (ret NodeSet) {
	return p.dataAtom(atom.P)
}

// Img returns NodeSet which node type is ElementNode and it's element type is `img`.
func (p NodeSet) Img() (ret NodeSet) {
	return p.dataAtom(atom.Img)
}

// Ol returns NodeSet which node type is ElementNode and it's element type is `ol`.
func (p NodeSet) Ol() (ret NodeSet) {
	return p.dataAtom(atom.Ol)
}

// Ul returns NodeSet which node type is ElementNode and it's element type is `ul`.
func (p NodeSet) Ul() (ret NodeSet) {
	return p.dataAtom(atom.Ul)
}

// Dl returns NodeSet which node type is ElementNode and it's element type is `dl`.
func (p NodeSet) Dl() (ret NodeSet) {
	return p.dataAtom(atom.Dl)
}

// Dt returns NodeSet which node type is ElementNode and it's element type is `dt`.
func (p NodeSet) Dt() (ret NodeSet) {
	return p.dataAtom(atom.Dt)
}

// Span returns NodeSet which node type is ElementNode and it's element type is `span`.
func (p NodeSet) Span() (ret NodeSet) {
	return p.dataAtom(atom.Span)
}

// Div returns NodeSet which node type is ElementNode and it's element type is `div`.
func (p NodeSet) Div() (ret NodeSet) {
	return p.dataAtom(atom.Div)
}

// Nav returns NodeSet which node type is ElementNode and it's element type is `nav`.
func (p NodeSet) Nav() (ret NodeSet) {
	return p.dataAtom(atom.Nav)
}

// Li returns NodeSet which node type is ElementNode and it's element type is `li`.
func (p NodeSet) Li() (ret NodeSet) {
	return p.dataAtom(atom.Li)
}

// Class returns NodeSet which `class` attribute is `v`.
func (p NodeSet) Class(v string) (ret NodeSet) {
	return p.Attribute__0("class", v)
}

// Id returns NodeSet which `id` attribute is `v`.
func (p NodeSet) Id(v string) (ret NodeSet) {
	return p.Attribute__0("id", v).One()
}

// -----------------------------------------------------------------------------

// ExactText returns text of NodeSet.
// exactlyOne=false: if NodeSet is more than one, returns first node's text (if
// node type is not TextNode, return error).
func (p NodeSet) ExactText__1(exactlyOne bool) (text string, err error) {
	node, err := p.CollectOne__1(exactlyOne)
	if err != nil {
		return
	}
	return exactText(node)
}

func (p NodeSet) ExactText__0() (text string, err error) {
	return p.ExactText__1(false)
}

// Text returns text of NodeSet.
// exactlyOne=false: if NodeSet is more than one, returns first node's text.
func (p NodeSet) Text__1(exactlyOne bool) (text string, err error) {
	node, err := p.CollectOne__1(exactlyOne)
	if err != nil {
		return
	}
	return textOf(node), nil
}

func (p NodeSet) Text__0() (text string, err error) {
	return p.Text__1(false)
}

// ScanInt returns int value of p.Text().
// exactlyOne=false: if NodeSet is more than one, returns first node's value.
func (p NodeSet) ScanInt(format string, exactlyOne ...bool) (v int, err error) {
	text, err := p.Text__1(exactlyOne != nil && exactlyOne[0])
	if err != nil {
		return
	}
	err = fmtSscanf(text, format, &v)
	if err != nil {
		v = 0
	}
	return
}

func fmtSscanf(text, format string, v *int) (err error) {
	prefix, suffix, err := parseFormat(format)
	if err != nil {
		return
	}
	if strings.HasPrefix(text, prefix) && strings.HasSuffix(text, suffix) {
		text = text[len(prefix) : len(text)-len(suffix)]
		*v, err = strconv.Atoi(strings.Replace(text, ",", "", -1))
		return
	}
	return ErrInvalidScanFormat
}

func parseFormat(format string) (prefix, suffix string, err error) {
	pos := strings.Index(format, "%d")
	if pos < 0 {
		pos = strings.Index(format, "%v")
	}
	if pos < 0 {
		err = ErrInvalidScanFormat
		return
	}
	prefix = strings.Replace(format[:pos], "%%", "%", -1)
	suffix = strings.Replace(format[pos+2:], "%%", "%", -1)
	return
}

// UnitedFloat returns UnitedFloat value of p.Text().
// exactlyOne=false: if NodeSet is more than one, returns first node's value.
func (p NodeSet) UnitedFloat__1(exactlyOne bool) (v float64, err error) {
	text, err := p.Text__1(exactlyOne)
	if err != nil {
		return
	}
	n := len(text)
	if n == 0 {
		return 0, ErrEmptyText
	}
	unit := 1.0
	switch text[n-1] {
	case 'k', 'K':
		unit = 1000
		text = text[:n-1]
	}
	v, err = strconv.ParseFloat(text, 64)
	if err != nil {
		return
	}
	return v * unit, nil
}

func (p NodeSet) UnitedFloat__0() (v float64, err error) {
	return p.UnitedFloat__1(false)
}

// Int returns int value of p.Text().
// exactlyOne=false: if NodeSet is more than one, returns first node's value.
func (p NodeSet) Int__1(exactlyOne bool) (v int, err error) {
	text, err := p.Text__1(exactlyOne)
	if err != nil {
		return
	}
	return strconv.Atoi(strings.Replace(text, ",", "", -1))
}

func (p NodeSet) Int__0() (v int, err error) {
	return p.Int__1(false)
}

// AttrVal returns attribute value of NodeSet.
// exactlyOne=false: if NodeSet is more than one, returns first node's attribute value.
func (p NodeSet) AttrVal(k string, exactlyOne ...bool) (text string, err error) {
	node, err := p.CollectOne__1(exactlyOne != nil && exactlyOne[0])
	if err != nil {
		return
	}
	return attributeVal(node, k)
}

// HrefVal returns href attribute's value of NodeSet.
// exactlyOne=false: if NodeSet is more than one, returns first node's attribute value.
func (p NodeSet) HrefVal__1(exactlyOne bool) (text string, err error) {
	return p.AttrVal("href", exactlyOne)
}

func (p NodeSet) HrefVal__0() (text string, err error) {
	return p.AttrVal("href", false)
}

// -----------------------------------------------------------------------------

// Href returns href attribute's value of NodeSet.
func (p NodeSet) Href__0() (text string, err error) {
	return p.AttrVal("href", false)
}

// Href returns NodeSet which `href` attribute is `v`.
func (p NodeSet) Href__1(v string) (ret NodeSet) {
	return p.Attribute__0("href", v)
}

// Href returns href attribute's value of NodeSet.
func (p NodeSet) Href__2(exactlyOne bool) (text string, err error) {
	return p.AttrVal("href", exactlyOne)
}

// HasAttr returns true if NodeSet has attribute k or not.
func (p NodeSet) HasAttr(k string, exactlyOne ...bool) bool {
	_, e := p.AttrVal(k, exactlyOne...)
	return e == nil
}

// Attr returns attribute value of NodeSet.
func (p NodeSet) Attr__0(k string, exactlyOne ...bool) (text string, err error) {
	return p.AttrVal(k, exactlyOne...)
}

func (p NodeSet) Attr__1(k, v string) (ret NodeSet) {
	return p.Attribute__0(k, v)
}

// -----------------------------------------------------------------------------
