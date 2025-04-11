package platform

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/redsift/go-siftjson"
)

// BlockState is a snapshot of the blocked sift versions and accounts.
type BlockState struct {
	accounts map[siftjson.Key]struct{}
	sifts    map[siftjson.ID]bool
}

// IsIdentityBlocked returns true if the sift version is blocked, and the hard-block flag.
func (b *BlockState) IsIdentityBlocked(id siftjson.ID) (blocked, hard bool) {
	hard, blocked = b.sifts[id]
	return
}

// IsInstanceKeyBlocked returns true if the account or sift version is blocked.
func (b *BlockState) IsInstanceKeyBlocked(key siftjson.InstanceKey) (blocked, hard bool) {
	if hard, blocked = b.sifts[key.ID]; blocked {
		return
	}
	_, blocked = b.accounts[key.Key]
	return
}

// IsKeyBlocked returns true if the account is blocked.
func (b *BlockState) IsKeyBlocked(key siftjson.Key) bool {
	_, blocked := b.accounts[key]
	return blocked
}

// BlockedService wraps BlockedAccount and BlockedSiftsSlice into an auto-updating service.
type BlockedService struct {
	initialized     atomic.Bool
	onError         func(error) bool
	blocked         atomic.Pointer[BlockState]
	blockedAccounts *BlockedAccountsSlice
	blockedSifts    *BlockedSiftsSlice
}

type BlockedServiceOption func(*BlockedService)

// WithErrorHandler registers an error handler. If the handler returns false, the service will shut
// down.
func WithErrorHandler(fn func(error) bool) BlockedServiceOption {
	return func(bs *BlockedService) {
		bs.onError = fn
	}
}

// NewBlockedService returns a new BlockedService with the given options. To start the auto-update
// process call Init.
func NewBlockedService(
	blockedAccounts *BlockedAccountsSlice,
	blockedSifts *BlockedSiftsSlice,
	opts ...BlockedServiceOption,
) *BlockedService {
	return &BlockedService{
		onError: func(error) bool { return true },
	}
}

// Load the current BlockState instance.
func (bs *BlockedService) Load() *BlockState {
	return bs.blocked.Load()
}

// LoadSiftVersions loads and returns the currently blocked accounts.
func (bs *BlockedService) LoadAccounts(ctx context.Context) ([]BlockedAccount, error) {
	return bs.blockedAccounts.Load(ctx)
}

// LoadSiftVersions loads and returns the currently blocked sift versions.
func (bs *BlockedService) LoadSiftVersions(ctx context.Context) ([]BlockedSiftVersion, error) {
	return bs.blockedSifts.Load(ctx)
}

// Init loads the current state and subscribes to the store in the backend.
func (bs *BlockedService) Init(ctx context.Context) error {
	bs.blocked.Store(&BlockState{
		accounts: map[siftjson.Key]struct{}{},
		sifts:    map[siftjson.ID]bool{},
	})

	if err := bs.blockedAccounts.Subscribe(ctx, func(updated []BlockedAccount, err error) bool {
		if err != nil {
			if !bs.onError(fmt.Errorf("received blocked account update error: %w", err)) {
				return false
			}
			return ctx.Err() == nil
		}
		bs.updateBlockedAccounts(updated)
		return ctx.Err() == nil
	}); err != nil {
		return err
	}

	if accounts, err := bs.blockedAccounts.Load(ctx); err != nil {
		return err
	} else {
		bs.updateBlockedAccounts(accounts)
	}

	if err := bs.blockedSifts.Subscribe(ctx, func(updated []BlockedSiftVersion, err error) bool {
		if err != nil {
			if !bs.onError(fmt.Errorf("received blocked sift versions update error: %w", err)) {
				return false
			}
			return ctx.Err() == nil
		}
		bs.updateBlockedSifts(updated)
		return true
	}); err != nil {
		return err
	}

	if sifts, err := bs.blockedSifts.Load(ctx); err != nil {
		return err
	} else {
		bs.updateBlockedSifts(sifts)
	}

	bs.initialized.Store(true)

	return nil
}

func (e *BlockedService) updateBlocked(fn func(*BlockState) *BlockState) {
	for {
		last := e.blocked.Load()
		next := fn(last)
		if e.blocked.CompareAndSwap(last, next) {
			return
		}
	}
}

func (e *BlockedService) updateBlockedAccounts(next []BlockedAccount) {
	accounts := map[siftjson.Key]struct{}{}
	for _, acc := range next {
		accounts[siftjson.Key{
			GUID:  acc.GUID,
			AccID: acc.Account,
		}] = struct{}{}
	}

	e.updateBlocked(func(last *BlockState) *BlockState {
		return &BlockState{
			accounts: accounts,
			sifts:    last.sifts,
		}
	})
}

func (e *BlockedService) updateBlockedSifts(next []BlockedSiftVersion) {
	sifts := map[siftjson.ID]bool{}
	for _, sv := range next {
		sifts[sv.ID] = sv.Hard
	}

	e.updateBlocked(func(last *BlockState) *BlockState {
		return &BlockState{
			accounts: last.accounts,
			sifts:    sifts,
		}
	})
}
