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

package hdq_test

import (
	"strings"
	"testing"

	"github.com/goplus/hdq"
	"github.com/goplus/hdq/fetcher/githubisstask"
	"github.com/goplus/hdq/fetcher/gopkg"
	"github.com/goplus/hdq/fetcher/torch"
	"github.com/goplus/hdq/hdqtest"

	repos "github.com/goplus/hdq/tutorial/02-GithubRepos"
)

func textOf(_ string, doc hdq.NodeSet) (ret string) {
	ret, _ = doc.Text__0()
	return
}

func TestText(t *testing.T) {
	hdqtest.FromDir(t, "", "./_testdata/text", textOf, "data.zip#index.htm", "zip")
}

func TestGithub(t *testing.T) {
	hdqtest.FromDir(t, "", "./_testdata/github", repos.New, "data.zip#index.htm", "zip")
}

func TestTorch(t *testing.T) {
	hdqtest.FromDir(t, "", "./fetcher/torch/_testdata", torch.New, "data.zip#index.htm", "zip")
}

func TestGoPkg(t *testing.T) {
	hdqtest.FromDir(t, "", "./fetcher/gopkg/_testdata", gopkg.New, "data.zip#index.htm", "zip")
}

func TestGithubIssueTask(t *testing.T) {
	hdqtest.FromDir(t, "", "./fetcher/githubisstask/_testdata", githubisstask.New, "data.zip#index.htm", "zip")
}

func TestSource(t *testing.T) {
	const data = "<html><body>hello</body></html>"
	doc := hdq.Source([]byte(data))
	sources := []any{
		[]byte(data),
		strings.NewReader(data),
		doc,
	}
	for _, in := range sources {
		v := hdq.Source(in)
		if text, err := v.Text__0(); err != nil || text != "hello" {
			t.Fatal("Source failed: ", text, err)
		}
	}
	if doc := hdq.Source("unknown:123"); doc.Ok() {
		t.Fatal("Source failed: no error?")
	}
	defer func() {
		if recover() == nil {
			t.Fatalf("Source failed: no panic?")
		}
	}()
	hdq.Source(123)
}

func TestErrNodeSet(t *testing.T) {
	docErr := hdq.NodeSet{Err: hdq.ErrInvalidNode}
	fns := []func(hdq.NodeSet) hdq.NodeSet{
		(hdq.NodeSet).Child,
		(hdq.NodeSet).Parent,
		(hdq.NodeSet).PrevSiblings,
		(hdq.NodeSet).NextSiblings,
		(hdq.NodeSet).Any,
	}
	for _, fn := range fns {
		if v := fn(docErr); v != docErr {
			t.Fatal("ErrNodeSet failed:", v)
		}
	}
	const data = "<html><body>hello</body></html>"
	doc := hdq.Source([]byte(data))
	for _, fn := range fns {
		fn(doc)
	}
}
