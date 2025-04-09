package platform

import (
	"context"
	"strings"

	"github.com/redsift/go-cfg/dcfg"
	"github.com/redsift/go-siftjson"
)

const KEY_ACCOUNTS = "accounts"

var BlockedAccountssV1Key = SliceKey[BlockedAccount](1, KEY_ACCOUNTS, KEY_BLOCKED)

type BlockedAccountsSlice = dcfg.TypedSlice[BlockedAccount]

type BlockedAccount struct {
	GUID    siftjson.GUID
	Account siftjson.AccountID
	Reason  string
}

func BlockedAccounts(b dcfg.Backend) *BlockedAccountsSlice {
	res, _ := dcfg.NewTypedSlice[BlockedAccount](b, BlockedSiftsV1Key, func(a, b BlockedAccount) int {
		if diff := strings.Compare(string(a.GUID), string(b.GUID)); diff != 0 {
			return diff
		}
		if diff := strings.Compare(string(a.Account), string(b.Account)); diff != 0 {
			return diff
		}
		return 0
	})
	return res
}

func LoadBlockedAccounts(ctx context.Context, b dcfg.Backend) (out []BlockedAccount, _ error) {
	return BlockedAccounts(b).Load(ctx)
}
