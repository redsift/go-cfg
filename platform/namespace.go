package platform

import "github.com/redsift/go-cfg/dcfg"

const (
	KEY_PLATFORM = "platform"
	KEY_SIFTS    = "sifts"
)

func Key(version dcfg.Version, key ...string) dcfg.Key {
	return dcfg.NewKey(version, KEY_PLATFORM, key...)
}
