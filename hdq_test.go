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
	"testing"

	"github.com/goplus/hdq"
	"github.com/goplus/hdq/hdqtest"
	"github.com/goplus/hdq/pysig/torch"
)

func textOf(doc hdq.NodeSet) (ret string) {
	ret, _ = doc.Text__0()
	return
}

func TestText(t *testing.T) {
	hdqtest.FromDir(t, "", "./_testdata/text", textOf)
}

func TestTestdata(t *testing.T) {
	hdqtest.FromDir(t, "", "./pysig/torch/_testdata", torch.New)
}
