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

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/goplus/hdq/cmd/hdq/internal/base"
	"github.com/goplus/hdq/cmd/hdq/internal/fetch"
	"github.com/goplus/hdq/cmd/hdq/internal/help"
	"github.com/qiniu/x/log"

	_ "github.com/goplus/hdq/fetcher/githubisstask"
	_ "github.com/goplus/hdq/fetcher/gopkg"
	_ "github.com/goplus/hdq/fetcher/hrefs"
	_ "github.com/goplus/hdq/fetcher/torch"
	_ "github.com/goplus/hdq/stream/http/cached"
)

func mainUsage() {
	help.PrintUsage(os.Stderr, base.Hdq)
	os.Exit(2)
}

func init() {
	flag.Usage = mainUsage
	base.Hdq.Commands = []*base.Command{
		fetch.Cmd,
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
	}
	log.SetFlags(log.Ldefault &^ log.LstdFlags)

	base.CmdName = args[0] // for error messages
	if args[0] == "help" {
		help.Help(os.Stderr, args[1:])
		return
	}

BigCmdLoop:
	for bigCmd := base.Hdq; ; {
		for _, cmd := range bigCmd.Commands {
			if cmd.Name() != args[0] {
				continue
			}
			args = args[1:]
			if len(cmd.Commands) > 0 {
				bigCmd = cmd
				if len(args) == 0 {
					help.PrintUsage(os.Stderr, bigCmd)
					os.Exit(2)
				}
				if args[0] == "help" {
					help.Help(os.Stderr, append(strings.Split(base.CmdName, " "), args[1:]...))
					return
				}
				base.CmdName += " " + args[0]
				continue BigCmdLoop
			}
			if !cmd.Runnable() {
				continue
			}
			cmd.Run(cmd, args)
			return
		}
		helpArg := ""
		if i := strings.LastIndex(base.CmdName, " "); i >= 0 {
			helpArg = " " + base.CmdName[:i]
		}
		fmt.Fprintf(os.Stderr, "hdq %s: unknown command\nRun 'hdq help%s' for usage.\n", base.CmdName, helpArg)
		os.Exit(2)
	}
}
