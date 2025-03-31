package platform

import (
	"context"

	"github.com/redsift/go-cfg/dcfg"
	"github.com/redsift/go-siftjson"
)

const KEY_ACCOUNTS = "accounts"

var BlockedAccountssV1Key = SliceKey[BlockedAccount](1, KEY_ACCOUNTS, KEY_BLOCKED)

type BlockedAccountsSlice = dcfg.TypedSlice[BlockedAccount]

type BlockedAccount struct {
	Account siftjson.AccountID
	Reason  string
}

func BlockedAccounts(b dcfg.Backend) *BlockedAccountsSlice {
	res, _ := dcfg.NewTypedSlice[BlockedAccount](b, BlockedSiftsV1Key)
	return res
}

func LoadBlockedAccounts(ctx context.Context, b dcfg.Backend) (out []BlockedAccount, _ error) {
	return BlockedAccounts(b).Load(ctx)
}
