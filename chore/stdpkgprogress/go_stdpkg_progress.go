package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/goplus/hdq/fetcher"
	"github.com/goplus/hdq/fetcher/githubisstask"
	_ "github.com/goplus/hdq/stream/http/nocache"
)

const importedBy = "Imported By: "

// Usage: stdpkgprogress
func main() {
	doc, err := fetcher.FromInput("githubisstask", "goplus/llgo#642")
	if err != nil {
		panic(err)
	}
	var done, total float64
	var ndone, ntotal int
	ret := doc.(githubisstask.Result)
	for _, task := range ret.Tasks {
		desc := task.Desc // fmt* (Imported By: 4513111)
		if pos := strings.Index(desc, "Imported By: "); pos > 0 {
			ntext := strings.TrimSuffix(desc[pos+len(importedBy):], ")")
			if n, e := strconv.Atoi(ntext); e == nil {
				w := math.Log2(float64(n) + 1)
				total += w
				ntotal += n
				if task.Done {
					done += w
					ndone += n
				}
			}
		}
	}
	np := float64(ndone) / float64(ntotal) * 100
	fmt.Printf("Progress: %.2f%% (%.2f%%)\n", done/total*100, np)
}
