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

package http

import (
	"errors"
	"io"
	"net/http"
)

var (
	// DefaultUserAgent is the default UserAgent and is used by HTTPSource.
	DefaultUserAgent string
	ReqHeaderProc    func(req *http.Request)
)

// -------------------------------------------------------------------------------------

// Open opens a http file object.
func Open(url string) (io.ReadCloser, error) {
	resp, err := Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if DefaultUserAgent != "" {
		req.Header.Set("User-Agent", DefaultUserAgent)
	}
	if ReqHeaderProc != nil {
		ReqHeaderProc(req)
	}
	if resp, err = http.DefaultClient.Do(req); err != nil {
		return
	}
	if resp.StatusCode/100 != 2 {
		resp.Body.Close()
		err = errors.New(resp.Status)
	}
	return
}

// -------------------------------------------------------------------------------------
