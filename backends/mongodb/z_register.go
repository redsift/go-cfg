package mongodb

import (
	"fmt"
	"net/url"

	"github.com/redsift/go-cfg/dcfg"
)

func init() {
	dcfg.Register("mongodb", func(uri *url.URL) (dcfg.Backend, error) {
		q := uri.Query()
		if !q.Has("db") {
			return nil, fmt.Errorf("missing mongo database name in %q", uri)
		}

		db := q.Get("db")
		q.Del("db")
		uri.RawQuery = q.Encode()

		return New(uri.String(), db)
	})
}
