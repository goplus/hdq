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

package torch

import (
	"testing"

	"github.com/goplus/hdq/hdqtest"
)

func TestTestdata(t *testing.T) {
	hdqtest.FromDir(t, "", "./_testdata", New)
}

func TestURL(t *testing.T) {
	if v := URL("eye"); v != "https://pytorch.org/docs/stable/generated/torch.eye.html" {
		t.Fatal(v)
	}
}
