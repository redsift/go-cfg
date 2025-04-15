package platform

import (
	"context"
	"testing"

	"github.com/redsift/go-cfg/backends/nilbe"
	"github.com/redsift/go-cfg/dcfg"
	"github.com/stretchr/testify/require"
)

func TestBlockedAccounts(t *testing.T) {
	blocked := BlockedAccounts(nilbe.Nil)
	require.NotNil(t, blocked)
	result, err := blocked.Load(context.TODO())
	require.Error(t, dcfg.ErrNotFound, err)
	require.Nil(t, result)
}
