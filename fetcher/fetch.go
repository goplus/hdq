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

package fetcher

import (
	"reflect"

	"github.com/goplus/hdq"
)

// func(url string, doc hdq.NodeSet) <any-object>
type Conv = any

// -----------------------------------------------------------------------------

// Convert converts a html source to an object.
func Convert(conv reflect.Value, url string, in any) any {
	doc := reflect.ValueOf(hdq.Source(in))
	out := conv.Call([]reflect.Value{reflect.ValueOf(url), doc})
	return out[0].Interface()
}

// -----------------------------------------------------------------------------

// New creates a new object from a html source by a registered converter.
func New(pageType string, url string, in any) any {
	page, ok := convs[pageType]
	if !ok {
		panic("fetcher: unknown pageType - " + pageType)
	}
	return Convert(page.Conv, url, in)
}

// FromInput creates a new object from the html source with the specified input name.
func FromInput(pageType string, input string) any {
	page, ok := convs[pageType]
	if !ok {
		panic("fetcher: unknown pageType - " + pageType)
	}
	url := page.URL(input)
	return Convert(page.Conv, url, url)
}

// sitePageType represents a site page type.
type sitePageType struct {
	Conv reflect.Value
	URL  func(string) string
}

var (
	convs = map[string]sitePageType{}
)

// Register registers a convType with a convert function.
func Register(pageType string, conv Conv, urlOf func(string) string) {
	vConv := reflect.ValueOf(conv)
	convs[pageType] = sitePageType{vConv, urlOf}
}

// -----------------------------------------------------------------------------
