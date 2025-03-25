package platform

import (
	"context"

	"github.com/redsift/go-cfg/dcfg"
	"github.com/redsift/go-siftjson"
)

const KEY_LABS = "labs"

var LabsSiftsV1Key = SliceKey[LabsSift](1, KEY_SIFTS, KEY_LABS)

type LabsSift struct {
	GUID  siftjson.GUID
	Class string
}

func LabsSifts(b dcfg.Backend) *dcfg.TypedSlice[LabsSift] {
	res, _ := dcfg.NewTypedSlice[LabsSift](b, BlockedSiftsV1Key)
	return res
}

func LoadLabsSifts(ctx context.Context, b dcfg.Backend) (out []LabsSift, _ error) {
	return LabsSifts(b).Load(ctx)
}
