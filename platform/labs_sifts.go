package platform

import (
	"context"
	"errors"

	"github.com/redsift/go-cfg/dcfg"
	"github.com/redsift/go-siftjson"
)

const KEY_LABS = "labs"

var LabsSiftsV1Key = MapKey[siftjson.GUID, string](1, KEY_SIFTS, KEY_LABS)

func LabsSifts(b dcfg.Backend) *dcfg.TypedMap[siftjson.GUID, string] {
	res, _ := dcfg.NewTypedMap[siftjson.GUID, string](b, LabsSiftsV1Key)
	return res
}

func LoadLabsSifts(ctx context.Context, b dcfg.Backend) (map[siftjson.GUID]string, error) {
	out, err := LabsSifts(b).Load(ctx)
	if errors.Is(err, dcfg.ErrNotFound) {
		return map[siftjson.GUID]string{}, nil
	}
	return out, err
}
