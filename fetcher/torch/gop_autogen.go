// Code generated by gop (Go+); DO NOT EDIT.

package torch

import (
	"github.com/goplus/hdq"
	"github.com/goplus/hdq/fetcher"
	"github.com/qiniu/x/errors"
	"strings"
)

const GopPackage = "github.com/goplus/hdq"
const _ = true
const spaces = " \t\r\n¶"

type Result struct {
	Name string `json:"name"`
	Doc  string `json:"doc"`
	Sig  string `json:"sig"`
	URL  string `json:"url,omitempty"`
}
//line fetcher/torch/pysig_torch.gop:38:1
// New creates a new Result from a html document.
func New(url string, doc hdq.NodeSet) Result {
//line fetcher/torch/pysig_torch.gop:40:1
	if doc.Ok() {
//line fetcher/torch/pysig_torch.gop:41:1
		fn := doc.Any().Dl().Class("py function")
//line fetcher/torch/pysig_torch.gop:42:1
		decl := func() (_gop_ret string) {
//line fetcher/torch/pysig_torch.gop:42:1
			var _gop_err error
//line fetcher/torch/pysig_torch.gop:42:1
			_gop_ret, _gop_err = fn.FirstElementChild().Dt().Text__0()
//line fetcher/torch/pysig_torch.gop:42:1
			if _gop_err != nil {
//line fetcher/torch/pysig_torch.gop:42:1
				_gop_err = errors.NewFrame(_gop_err, "fn.firstElementChild.dt.text", "fetcher/torch/pysig_torch.gop", 42, "torch.New")
//line fetcher/torch/pysig_torch.gop:42:1
				panic(_gop_err)
			}
//line fetcher/torch/pysig_torch.gop:42:1
			return
		}()
//line fetcher/torch/pysig_torch.gop:43:1
		pos := strings.IndexByte(decl, '(')
//line fetcher/torch/pysig_torch.gop:44:1
		if pos > 0 {
//line fetcher/torch/pysig_torch.gop:45:1
			name := strings.TrimPrefix(decl[:pos], "torch.")
//line fetcher/torch/pysig_torch.gop:46:1
			sig := decl[pos:]
//line fetcher/torch/pysig_torch.gop:47:1
			return Result{strings.TrimSpace(name), "", strings.TrimRight(sig, spaces), url}
		}
	}
//line fetcher/torch/pysig_torch.gop:50:1
	return Result{"", "", "<NULL>", url}
}
//line fetcher/torch/pysig_torch.gop:53:1
// URL returns the input URL for the given name.
func URL(name interface{}) string {
//line fetcher/torch/pysig_torch.gop:55:1
	return "https://pytorch.org/docs/stable/generated/torch." + name.(string) + ".html"
}
//line fetcher/torch/pysig_torch.gop:58:1
func init() {
//line fetcher/torch/pysig_torch.gop:59:1
	fetcher.Register("torch", New, URL)
}
