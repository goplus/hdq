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
	"bytes"
	"errors"
	"io"

	"github.com/goplus/hdq/stream"
	"golang.org/x/net/html"

	_ "github.com/goplus/hdq/stream/zip"
)

const (
	GopPackage = true // to indicate this is a Go+ package
)

var (
	ErrNotFound = errors.New("entity not found")
	ErrBreak    = errors.New("break")

	ErrTooManyNodes = errors.New("too many nodes")
	ErrInvalidNode  = errors.New("invalid node")

	// ErrEmptyText represents an `empty text` error.
	ErrEmptyText = errors.New("empty text")

	// ErrInvalidScanFormat represents an `invalid fmt.Scan format` error.
	ErrInvalidScanFormat = errors.New("invalid fmt.Scan format")
)

// -----------------------------------------------------------------------------

type NodeEnum interface {
	ForEach(filter func(node *html.Node) error)
}

type cachedGetter interface {
	Cached() int
}

// NodeSet represents a set of nodes.
type NodeSet struct {
	Data NodeEnum
	Err  error
}

// New creates a NodeSet object.
func New(r io.Reader) NodeSet {
	doc, err := html.Parse(r)
	if err != nil {
		return NodeSet{Err: err}
	}
	return NodeSet{Data: oneNode{doc}}
}

// Source opens a stream (if necessary) to create a NodeSet object.
func Source(r interface{}) (ret NodeSet) {
	switch v := r.(type) {
	case string:
		f, err := stream.Open(v)
		if err != nil {
			return NodeSet{Err: err}
		}
		return New(f)
	case []byte:
		r := bytes.NewReader(v)
		return New(r)
	case io.Reader:
		return New(v)
	case NodeSet: // input is a node set
		return v
	default:
		panic("unsupport source type")
	}
}

func (p NodeSet) Ok() bool {
	return p.Err == nil
}

func (p NodeSet) All() NodeSet {
	if _, ok := p.Data.(cachedGetter); ok {
		return p
	}
	nodes, err := p.Collect()
	if err != nil {
		return NodeSet{Err: err}
	}
	return NodeSet{Data: &fixNodes{nodes}}
}

func (p NodeSet) Gop_Enum(callback func(node NodeSet)) {
	if p.Err == nil {
		p.Data.ForEach(func(node *html.Node) error {
			t := NodeSet{Data: oneNode{node}}
			callback(t)
			return nil
		})
	}
}

func (p NodeSet) ForEach(callback func(node NodeSet)) {
	p.Gop_Enum(callback)
}

// Render renders the node set to the given writer.
func (p NodeSet) Render(w io.Writer, suffix ...string) (err error) {
	if p.Err != nil {
		return p.Err
	}
	p.Data.ForEach(func(node *html.Node) error {
		if e := html.Render(w, node); e != nil {
			err = e
			return ErrBreak
		}
		if suffix != nil {
			io.WriteString(w, suffix[0])
		}
		return nil
	})
	return
}

// -----------------------------------------------------------------------------

type oneNode struct {
	*html.Node
}

func (p oneNode) ForEach(filter func(node *html.Node) error) {
	filter(p.Node)
}

func (p oneNode) Cached() int {
	return 1
}

// -----------------------------------------------------------------------------

type fixNodes struct {
	nodes []*html.Node
}

func (p *fixNodes) ForEach(filter func(node *html.Node) error) {
	for _, node := range p.nodes {
		if filter(node) == ErrBreak {
			return
		}
	}
}

func (p *fixNodes) Cached() int {
	return len(p.nodes)
}

// Nodes creates a node set from the given nodes.
func Nodes(nodes ...*html.Node) (ret NodeSet) {
	return NodeSet{Data: &fixNodes{nodes}}
}

// -----------------------------------------------------------------------------

const (
	unknownNumNodes = -1
)

type anyNodes struct {
	data NodeEnum
}

func (p *anyNodes) ForEach(filter func(node *html.Node) error) {
	p.data.ForEach(func(node *html.Node) error {
		anyForEach(node, filter)
		return nil
	})
}

func (p *anyNodes) Cached() int {
	return unknownNumNodes
}

func anyForEach(p *html.Node, filter func(node *html.Node) error) error {
	if err := filter(p); err == nil || err == ErrBreak {
		return err
	}
	for node := p.FirstChild; node != nil; node = node.NextSibling {
		if anyForEach(node, filter) == ErrBreak {
			return ErrBreak
		}
	}
	return nil
}

// Any returns the all nodes as a node set.
func (p NodeSet) Any() (ret NodeSet) {
	if p.Err != nil {
		return p
	}
	return NodeSet{Data: &anyNodes{p.Data}}
}

// -----------------------------------------------------------------------------

type childLevelNodes struct {
	data  NodeEnum
	level int
}

func (p *childLevelNodes) ForEach(filter func(node *html.Node) error) {
	p.data.ForEach(func(node *html.Node) error {
		return childLevelForEach(node, p.level, filter)
	})
}

func childLevelForEach(p *html.Node, level int, filter func(node *html.Node) error) error {
	if level == 0 {
		return filter(p)
	}
	level--
	for node := p.FirstChild; node != nil; node = node.NextSibling {
		if childLevelForEach(node, level, filter) == ErrBreak {
			return ErrBreak
		}
	}
	return nil
}

type parentLevelNodes struct {
	data  NodeEnum
	level int
}

func (p *parentLevelNodes) ForEach(filter func(node *html.Node) error) {
	p.data.ForEach(func(node *html.Node) error {
		return parentLevelForEach(node, p.level, filter)
	})
}

func parentLevelForEach(p *html.Node, level int, filter func(node *html.Node) error) error {
	for level < 0 {
		if p = p.Parent; p == nil {
			return ErrNotFound
		}
		level++
	}
	return filter(p)
}

// Child returns the child node set. It is equivalent to ChildN(1).
func (p NodeSet) Child() (ret NodeSet) {
	return p.ChildN(1)
}

// ChildN returns the child node set at the given level.
func (p NodeSet) ChildN(level int) (ret NodeSet) {
	if p.Err != nil || level == 0 {
		return p
	}
	if level > 0 {
		return NodeSet{Data: &childLevelNodes{p.Data, level}}
	}
	return NodeSet{Data: &parentLevelNodes{p.Data, level}}
}

// Parent returns the parent node set. It is equivalent to ParentN(1).
func (p NodeSet) Parent() (ret NodeSet) {
	return p.ChildN(-1)
}

// ParentN returns the parent node set at the given level.
func (p NodeSet) ParentN(level int) (ret NodeSet) {
	return p.ChildN(-level)
}

// One returns the first node as a node set.
func (p NodeSet) One() (ret NodeSet) {
	if _, ok := p.Data.(oneNode); ok {
		return p
	}
	node, err := p.CollectOne__1(false)
	if err != nil {
		return NodeSet{Err: err}
	}
	return NodeSet{Data: oneNode{node}}
}

// -----------------------------------------------------------------------------

type siblingNodes struct {
	data  NodeEnum
	delta int
}

func (p *siblingNodes) ForEach(filter func(node *html.Node) error) {
	p.data.ForEach(func(node *html.Node) error {
		return siblingForEach(node, p.delta, filter)
	})
}

func siblingForEach(p *html.Node, delta int, filter func(node *html.Node) error) error {
	for delta > 0 {
		if p = p.NextSibling; p == nil {
			return ErrNotFound
		}
		delta--
	}
	for delta < 0 {
		if p = p.PrevSibling; p == nil {
			return ErrNotFound
		}
		delta++
	}
	return filter(p)
}

func (p NodeSet) NextSibling(delta int) (ret NodeSet) {
	if p.Err != nil {
		return p
	}
	return NodeSet{Data: &siblingNodes{p.Data, delta}}
}

func (p NodeSet) PrevSibling(delta int) (ret NodeSet) {
	return p.NextSibling(-delta)
}

// -----------------------------------------------------------------------------

type prevSiblingNodes struct {
	data NodeEnum
}

func (p *prevSiblingNodes) ForEach(filter func(node *html.Node) error) {
	p.data.ForEach(func(node *html.Node) error {
		for p := node.PrevSibling; p != nil; p = p.PrevSibling {
			if filter(p) == ErrBreak {
				return ErrBreak
			}
		}
		return nil
	})
}

func (p NodeSet) PrevSiblings() (ret NodeSet) {
	if p.Err != nil {
		return p
	}
	return NodeSet{Data: &prevSiblingNodes{p.Data}}
}

// -----------------------------------------------------------------------------

type nextSiblingNodes struct {
	data NodeEnum
}

func (p *nextSiblingNodes) ForEach(filter func(node *html.Node) error) {
	p.data.ForEach(func(node *html.Node) error {
		for p := node.NextSibling; p != nil; p = p.NextSibling {
			if filter(p) == ErrBreak {
				return ErrBreak
			}
		}
		return nil
	})
}

func (p NodeSet) NextSiblings() (ret NodeSet) {
	if p.Err != nil {
		return p
	}
	return NodeSet{Data: &nextSiblingNodes{p.Data}}
}

// -----------------------------------------------------------------------------

type firstChildNodes struct {
	data     NodeEnum
	nodeType html.NodeType
}

func (p *firstChildNodes) ForEach(filter func(node *html.Node) error) {
	p.data.ForEach(func(node *html.Node) error {
		child, err := firstChild(node, p.nodeType)
		if err != nil {
			return err
		}
		return filter(child)
	})
}

func (p NodeSet) FirstChild(nodeType html.NodeType) (ret NodeSet) {
	if p.Err != nil {
		return p
	}
	return NodeSet{Data: &firstChildNodes{p.Data, nodeType}}
}

func (p NodeSet) FirstTextChild() (ret NodeSet) {
	return p.FirstChild(html.TextNode)
}

func (p NodeSet) FirstElementChild() (ret NodeSet) {
	return p.FirstChild(html.ElementNode)
}

// -----------------------------------------------------------------------------

type lastChildNodes struct {
	data     NodeEnum
	nodeType html.NodeType
}

func (p *lastChildNodes) ForEach(filter func(node *html.Node) error) {
	p.data.ForEach(func(node *html.Node) error {
		child, err := lastChild(node, p.nodeType)
		if err != nil {
			return err
		}
		return filter(child)
	})
}

func (p NodeSet) LastChild(nodeType html.NodeType) (ret NodeSet) {
	if p.Err != nil {
		return p
	}
	return NodeSet{Data: &lastChildNodes{p.Data, nodeType}}
}

func (p NodeSet) LastTextChild() (ret NodeSet) {
	return p.LastChild(html.TextNode)
}

func (p NodeSet) LastElementChild() (ret NodeSet) {
	return p.LastChild(html.ElementNode)
}

// -----------------------------------------------------------------------------

type matchedNodes struct {
	data   NodeEnum
	filter func(node *html.Node) bool
}

func (p *matchedNodes) ForEach(filter func(node *html.Node) error) {
	p.data.ForEach(func(node *html.Node) error {
		if p.filter(node) {
			return filter(node)
		}
		return ErrNotFound
	})
}

// Match returns the matched node set.
func (p NodeSet) Match(filter func(node *html.Node) bool) (ret NodeSet) {
	if p.Err != nil {
		return p
	}
	return NodeSet{Data: &matchedNodes{p.Data, filter}}
}

// -----------------------------------------------------------------------------

type textNodes struct {
	data      NodeEnum
	doReplace bool
}

func (p *textNodes) ForEach(filter func(node *html.Node) error) {
	p.data.ForEach(func(t *html.Node) error {
		node := &html.Node{
			Parent: t,
			Type:   html.TextNode,
			Data:   textOf(t),
		}
		if p.doReplace {
			t.FirstChild = node
			t.LastChild = node
		}
		return filter(node)
	})
}

func (p NodeSet) ChildrenAsText(doReplace bool) (ret NodeSet) {
	if p.Err != nil {
		return p
	}
	return NodeSet{Data: &textNodes{p.Data, doReplace}}
}

// -----------------------------------------------------------------------------

// CollectOne returns the first node.
// If `exactly` is true, it will return an error if there are more than one node.
func (p NodeSet) CollectOne__1(exactly bool) (item *html.Node, err error) {
	if p.Err != nil {
		return nil, p.Err
	}
	err = ErrNotFound
	if exactly {
		p.Data.ForEach(func(node *html.Node) error {
			if err == ErrNotFound {
				item, err = node, nil
				return nil
			}
			err = ErrTooManyNodes
			return ErrBreak
		})
	} else {
		p.Data.ForEach(func(node *html.Node) error {
			item, err = node, nil
			return ErrBreak
		})
	}
	return
}

// CollectOne returns the first node.
func (p NodeSet) CollectOne__0() (item *html.Node, err error) {
	return p.CollectOne__1(false)
}

// Collect returns all nodes.
func (p NodeSet) Collect() (items []*html.Node, err error) {
	if p.Err != nil {
		return nil, p.Err
	}
	p.Data.ForEach(func(node *html.Node) error {
		items = append(items, node)
		return nil
	})
	return
}

// -----------------------------------------------------------------------------
