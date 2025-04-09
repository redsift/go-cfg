package nilbe

import (
	"net/url"

	"github.com/redsift/go-cfg/dcfg"
)

func init() {
	dcfg.Register("nil", func(*url.URL) (dcfg.Backend, error) {
		return Nil, nil
	})
}
