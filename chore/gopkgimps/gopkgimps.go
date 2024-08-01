/*
Copyright 2024 The GoPlus Authors (goplus.org)
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

package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/goplus/hdq/fetcher"
	"github.com/goplus/hdq/fetcher/gopkg"
	_ "github.com/goplus/hdq/stream/http/cached"
)

// Usage: gopkgimps [pkgPath ...]
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: gopkgimps [pkgPath ...]")
		os.Exit(1)
	}
	names := os.Args[1:]
	docs := make([]gopkg.Result, 0, len(names))
	for _, name := range names {
		log.Println("==> Fetch", name)
		doc, err := fetcher.FromInput("gopkg", name)
		if err != nil {
			panic(err)
		}
		docs = append(docs, doc.(gopkg.Result))
	}
	sort.Slice(docs, func(i, j int) bool {
		return docs[i].ImportedBy > docs[j].ImportedBy
	})
	for _, doc := range docs {
		if doc.ImportedBy == 0 {
			break
		}
		fmt.Printf("- [ ] %s (Imported By: %d)\n", doc.Path, doc.ImportedBy)
	}
}
