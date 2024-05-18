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

package stream

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"strings"
)

var (
	ErrUnknownScheme = errors.New("unknown scheme")
)

// -------------------------------------------------------------------------------------

type OpenFunc = func(file string) (io.ReadCloser, error)

var (
	openers = map[string]OpenFunc{}
)

// Register registers a scheme with an open function.
func Register(scheme string, open OpenFunc) {
	openers[scheme] = open
}

func Open(url string) (io.ReadCloser, error) {
	scheme := schemeOf(url)
	if scheme == "" {
		return os.Open(url)
	}
	if open, ok := openers[scheme]; ok {
		return open(url)
	}
	return nil, &fs.PathError{Op: "hdq/stream.Open", Err: ErrUnknownScheme, Path: url}
}

func schemeOf(url string) (scheme string) {
	pos := strings.IndexAny(url, ":/")
	if pos > 0 {
		if url[pos] == ':' {
			return url[:pos]
		}
	}
	return ""
}

// -------------------------------------------------------------------------------------
