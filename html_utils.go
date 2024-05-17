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
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// -----------------------------------------------------------------------------

// containsClass returns true if v is a class in source.
// eg. `ContainsClass("top current", "current")` returns true.
func containsClass(source string, v string) bool {
	for {
		pos := strings.IndexByte(source, ' ')
		if pos < 0 {
			return source == v
		}
		if source[:pos] == v {
			return true
		}
		source = source[pos+1:]
	}
}

// attributeVal returns value of the attribute `k`.
func attributeVal(node *html.Node, k string) (v string, err error) {
	if node.Type != html.ElementNode {
		return "", ErrInvalidNode
	}
	for _, attr := range node.Attr {
		if attr.Key == k {
			return attr.Val, nil
		}
	}
	return "", ErrNotFound
}

// firstChild returns the first child with type `nodeType` of the node `node`.
func firstChild(node *html.Node, nodeType html.NodeType) (p *html.Node, err error) {
	for p = node.FirstChild; p != nil; p = p.NextSibling {
		if p.Type == nodeType {
			return p, nil
		}
	}
	return nil, ErrNotFound
}

// lastChild returns the last child with type `nodeType` of the node `node`.
func lastChild(node *html.Node, nodeType html.NodeType) (p *html.Node, err error) {
	for p = node.LastChild; p != nil; p = p.PrevSibling {
		if p.Type == nodeType {
			return p, nil
		}
	}
	return nil, ErrNotFound
}

// -----------------------------------------------------------------------------

const (
	spaces = " \t\r\n"
)

// childEqualText returns true if the type of node's child is TextNode and it's Data equals `text`.
func childEqualText(node *html.Node, text string) bool {
	p := node.FirstChild
	if p == nil || p.NextSibling != nil {
		return false
	}
	return equalText(p, text)
}

// equalText returns true if the type of node is TextNode and it's Data equals `text`.
func equalText(node *html.Node, text string) bool {
	if node.Type != html.TextNode {
		return false
	}
	return node.Data == text
}

// containsText returns true if the type of node is TextNode and it's Data contains `text`.
func containsText(node *html.Node, text string) bool {
	if node.Type != html.TextNode {
		return false
	}
	return strings.Contains(node.Data, text)
}

// hasPrefixText returns true if the type of node is TextNode and its Data has prefix `text`.
func hasPrefixText(node *html.Node, text string) bool {
	if node.Type != html.TextNode {
		return false
	}
	return strings.Contains(strings.TrimLeft(node.Data, spaces), text)
}

// exactText returns text of node if the type of node is TextNode.
func exactText(node *html.Node) (string, error) {
	if node.Type != html.TextNode {
		return node.Data, nil
	}
	return "", ErrInvalidNode
}

// textOf returns text data of node's all childs.
func textOf(node *html.Node) string {
	var printer textPrinter
	printer.printNode(node)
	return string(printer.data)
}

type textPrinter struct {
	data         []byte
	notLineStart bool
	hasSpace     bool
}

func (p *textPrinter) printText(v string, hasRightSpace bool) {
	if v == "" {
		return
	}
	if p.notLineStart && p.hasSpace {
		p.data = append(p.data, ' ')
	} else {
		p.notLineStart = true
	}
	p.data = append(p.data, v...)
	p.hasSpace = hasRightSpace
}

func (p *textPrinter) printNode(node *html.Node) {
	if node == nil {
		return
	}
	if node.Type == html.TextNode {
		p.printText(textTrimRight(textTrimLeft(node.Data, &p.hasSpace)))
		return
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		p.printNode(child)
	}
	switch node.DataAtom {
	case atom.P:
		p.data = append(p.data, '\n')
		p.notLineStart = false
	}
}

func textTrimLeft(v string, hasSpace *bool) string {
	ret := strings.TrimLeft(v, spaces)
	if len(v) != len(ret) {
		*hasSpace = true
	}
	return ret
}

func textTrimRight(v string) (string, bool) {
	ret := strings.TrimRight(v, spaces)
	return ret, len(v) != len(ret)
}

// -----------------------------------------------------------------------------
