hdq - HTML DOM Query Language for Go+
========

[![Build Status](https://github.com/goplus/hdq/actions/workflows/go.yml/badge.svg)](https://github.com/goplus/hdq/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/goplus/hdq)](https://goreportcard.com/report/github.com/goplus/hdq)
[![GitHub release](https://img.shields.io/github/v/tag/goplus/hdq.svg?label=release)](https://github.com/goplus/hdq/releases)
[![Coverage Status](https://codecov.io/gh/goplus/hdq/branch/main/graph/badge.svg)](https://codecov.io/gh/goplus/hdq)
[![Language](https://img.shields.io/badge/language-Go+-blue.svg)](https://github.com/goplus/gop)
[![GoDoc](https://img.shields.io/badge/godoc-reference-teal.svg)](https://pkg.go.dev/mod/github.com/goplus/hdq)

## Summary about hdq

hdq is a Go+ package for processing HTML documents.

## Tutorials

### Collect links of a html page

How to collect all links of a html page? If you use `hdq`, it is very easy.

```go
import "github.com/goplus/hdq"

func links(url any) []string {
	doc := hdq.Source(url)
	return [link for a <- doc.any.a if link := a.href?:""; link != ""]
}
```

At first, we call `hdq.Source(url)` to create a `node set` named `doc`. `doc` is a node set which only contains one node, the root node.

Then, select all `a` elements by `doc.any.a`. Here `doc.any` means all nodes in the html document.

Then, we visit all these `a` elements, get `href` attribute value and assign it to the variable `link`. If link is not empty, collect it.

At last, we return all collected links. Goto [tutorial/01-Links](tutorial/01-Links/links.gop) to get the full source code.
