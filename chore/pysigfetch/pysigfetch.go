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
	"os"

	"github.com/goplus/hdq/fetcher"
	_ "github.com/goplus/hdq/fetcher/torch"
)

type module struct {
	Name  string `json:"name"`
	Items []any  `json:"items"`
}

// Usage: pysigfetch pageType [name ...]
func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: pysigfetch pageType [name ...]")
		os.Exit(1)
	}
	pageType := os.Args[1]
	names := os.Args[2:]
	docs := make([]any, len(names))
	for i, name := range names {
		docs[i] = fetcher.FromInput(pageType, name)
	}
	json.NewEncoder(os.Stdout).Encode(module{pageType, docs})
}