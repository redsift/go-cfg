package platform

import (
	"context"

	"github.com/redsift/go-cfg/dcfg"
	"github.com/redsift/go-siftjson"
)

const KEY_BLOCKED = "blocked"

var BlockedSiftsV1Key = SliceKey[BlockedSiftVersion](1, KEY_SIFTS, KEY_BLOCKED)

type BlockedSiftsSlice = dcfg.TypedSlice[BlockedSiftVersion]

type BlockedSiftVersion struct {
	GUID   siftjson.GUID
	ID     siftjson.ID
	Reason string
}

func BlockedSifts(b dcfg.Backend) *BlockedSiftsSlice {
	res, _ := dcfg.NewTypedSlice[BlockedSiftVersion](b, BlockedSiftsV1Key)
	return res
}

func LoadBlockedSifts(ctx context.Context, b dcfg.Backend) (out []BlockedSiftVersion, _ error) {
	return BlockedSifts(b).Load(ctx)
}
