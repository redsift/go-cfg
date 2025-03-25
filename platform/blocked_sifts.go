package platform

import (
	"context"

	"github.com/redsift/go-cfg/dcfg"
	"github.com/redsift/go-siftjson"
)

const KEY_BLOCKED = "blocked"

var BlockedSiftsV1Key = Key(1, KEY_SIFTS, KEY_BLOCKED)

type BlockedSiftVersion struct {
	GUID   siftjson.GUID
	ID     siftjson.ID
	Reason string
}

func BlockedSifts(b dcfg.Backend) *dcfg.TypedSlice[BlockedSiftVersion] {
	return dcfg.NewTypedSlice[BlockedSiftVersion](b, BlockedSiftsV1Key)
}

func LoadBlockedSifts(ctx context.Context, b dcfg.Backend) (out []BlockedSiftVersion, _ error) {
	return BlockedSifts(b).Load(ctx)
}
