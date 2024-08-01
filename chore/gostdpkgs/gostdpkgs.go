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
	"os"
	"runtime"
	"strings"
)

func main() {
	dir := runtime.GOROOT() + "/src/"
	fis, err := os.ReadDir(dir)
	check(err)
	pkgs := collect(nil, fis, dir, "")
	fmt.Println(strings.Join(pkgs, "\n"))
}

func collect(pkgs []string, fis []os.DirEntry, dir, base string) []string {
	for _, fi := range fis {
		if !fi.IsDir() {
			continue
		}
		if name := fi.Name(); name != "cmd" && name != "internal" && name != "vendor" && name != "testdata" {
			nameSlash := name + "/"
			pkgDir := dir + nameSlash
			pkgFis, err := os.ReadDir(pkgDir)
			check(err)
			if hasGoFiles(pkgFis) {
				pkgs = append(pkgs, base+name)
			}
			pkgs = collect(pkgs, pkgFis, pkgDir, base+nameSlash)
		}
	}
	return pkgs
}

func hasGoFiles(fis []os.DirEntry) bool {
	for _, fi := range fis {
		if !fi.IsDir() && strings.HasSuffix(fi.Name(), ".go") {
			return true
		}
	}
	return false
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
