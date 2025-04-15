package mongodb

import (
	"testing"

	"github.com/redsift/go-cfg/platform"
	"github.com/stretchr/testify/require"
)

func TestCollectionName(t *testing.T) {
	be := &Backend{}

	require.Equal(t, "platform_v1_go_cfg_platform_BlockedAccount_slice", be.collectionName(platform.BlockedAccountssV1Key))
	require.Equal(t, "platform_v1_go_cfg_platform_BlockedSiftVersion_slice", be.collectionName(platform.BlockedSiftsV1Key))
}
