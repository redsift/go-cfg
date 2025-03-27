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
	impls[scheme] = factory
}
