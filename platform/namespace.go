package platform

import "github.com/redsift/go-cfg/dcfg"

const (
	KEY_PLATFORM = "platform"
	KEY_SIFTS    = "sifts"
)

func MapKey[K ~string, V any](version dcfg.Version, key ...string) dcfg.Key {
	return dcfg.NewKey[map[K]V](version, KEY_PLATFORM, key...)
}

func SliceKey[T any](version dcfg.Version, key ...string) dcfg.Key {
	return dcfg.NewKey[[]T](version, KEY_PLATFORM, key...)
}
