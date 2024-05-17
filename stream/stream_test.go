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
	"io"
	"strings"
	"testing"
)

type nilCloser struct {
	io.Reader
}

func (p *nilCloser) Close() error {
	return nil
}

func inlOpen(file string) (io.ReadCloser, error) {
	r := strings.NewReader(file)
	return &nilCloser{r}, nil
}

func TestBasic(t *testing.T) {
	RegisterSchema("inl", inlOpen)
	f, err := Open("inl://hello")
	if err != nil {
		t.Fatal("Open failed:", err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatal("ioutil.ReadAll failed:", err)
	}
	if string(b) != "hello" {
		t.Fatal("unexpected data")
	}
}

func TestUnknownSchema(t *testing.T) {
	_, err := Open("bad://foo")
	if err == nil || err.Error() != "stream.Open: unsupported schema - bad" {
		t.Fatal("Open failed:", err)
	}
}

func TestOpenFile(t *testing.T) {
	_, err := Open("/bin/not-exists/foo")
	if err == nil {
		t.Fatal("Open local file success?")
	}
}
