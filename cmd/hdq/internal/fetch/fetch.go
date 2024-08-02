/*
 * Copyright (c) 2024 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package fetch implements the "hdq fetch" command.
package fetch

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"github.com/goplus/hdq/cmd/hdq/internal/base"
	"github.com/goplus/hdq/fetcher"
)

// hdq fetch
var Cmd = &base.Command{
	UsageLine: "hdq fetch [flags] pageType [input ...]",
	Short:     "Fetch objects from the html source with the specified pageType and input",
}

func init() {
	Cmd.Run = runCmd
}

func runCmd(cmd *base.Command, args []string) {
	if len(args) < 2 {
		cmd.Usage(os.Stderr)
		return
	}
	pageType := args[0]
	inputs := args[1:]
	if len(inputs) == 1 && inputs[0] == "-" {
		b, _ := io.ReadAll(os.Stdin)
		inputs = strings.Split(strings.TrimSpace(string(b)), " ")
	}
	docs := make([]any, 0, len(inputs))
	for _, input := range inputs {
		log.Println("==> Fetch", input)
		doc, err := fetcher.FromInput(pageType, input)
		if err != nil {
			panic(err)
		}
		docs = append(docs, doc)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(docs)
}
