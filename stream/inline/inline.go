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

package inline

import (
	"io"
	"strings"

	"github.com/goplus/hdq/stream"
)

type nilCloser struct {
	io.Reader
}

func (p *nilCloser) Close() error {
	return nil
}

// Open opens a inline text object.
func Open(url string) (io.ReadCloser, error) {
	file := strings.TrimPrefix(url, "inline:")
	r := strings.NewReader(file)
	return &nilCloser{r}, nil
}

func init() {
	stream.Register("inline", Open)
}
