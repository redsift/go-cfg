package dcfg

import "net/url"

var (
	impls map[string]func(*url.URL) (Backend, error)
)

func Register(scheme string, factory func(*url.URL) (Backend, error)) {
	_, existing := impls[scheme]
	if existing {
		panic("two dcfg.Backend implementation for scheme \"+scheme+\"")
	}
	if impls == nil {
		impls = make(map[string]func(*url.URL) (Backend, error))
	}
	impls[scheme] = factory
}
