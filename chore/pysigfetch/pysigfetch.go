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
	_ "github.com/goplus/hdq/fetcher/torch"
	_ "github.com/goplus/hdq/stream/http/cached"
)

type module struct {
	Name  string `json:"name"`
	Items []any  `json:"items"`
}

// Usage: pysigfetch module [name ...]
func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: pysigfetch module [name ...]")
		os.Exit(1)
	}
	moduleName := os.Args[1]
	names := os.Args[2:]
	if len(names) == 1 && names[0] == "-" {
		b, _ := io.ReadAll(os.Stdin)
		names = strings.Split(strings.TrimSpace(string(b)), " ")
	}
	docs := make([]any, 0, len(names))
	for _, name := range names {
		log.Println("==> Fetch", name)
		doc, err := fetcher.FromInput(moduleName, name)
		if err != nil {
			panic(err)
		}
		docs = append(docs, doc)
	}
	json.NewEncoder(os.Stdout).Encode(module{moduleName, docs})
}
