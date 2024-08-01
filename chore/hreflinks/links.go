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
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/goplus/hdq/fetcher"
	_ "github.com/goplus/hdq/fetcher/hrefs"
	_ "github.com/goplus/hdq/stream/http/nocache"
)

// Usage: hreflinks [url ...]
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: hreflinks [url ...]")
		os.Exit(1)
	}
	urls := os.Args[1:]
	if len(urls) == 1 && urls[0] == "-" {
		b, _ := io.ReadAll(os.Stdin)
		urls = strings.Split(strings.TrimSpace(string(b)), "\n")
	}
	docs := make([]any, 0, len(urls))
	for _, url := range urls {
		log.Println("==> Fetch", url)
		doc, err := fetcher.FromInput("hrefs", url)
		if err != nil {
			panic(err)
		}
		docs = append(docs, doc)
	}
	json.NewEncoder(os.Stdout).Encode(docs)
}
