package platform

import "github.com/redsift/go-cfg/dcfg"

const (
	KEY_PLATFORM = "platform"
	KEY_SIFTS    = "sifts"
)

func SliceKey[T any](version dcfg.Version, key ...string) dcfg.Key {
	return dcfg.NewKey[[]T](version, KEY_PLATFORM, key...)
}
