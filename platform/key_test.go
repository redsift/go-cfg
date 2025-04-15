package platform_test

import (
	"testing"

	"github.com/redsift/go-cfg/platform"
	"github.com/stretchr/testify/require"
)

func TestKeyStrings(t *testing.T) {
	require.Equal(t, "platform/v1/accounts/blocked:go-cfg/platform/BlockedAccount_slice", platform.BlockedAccountssV1Key.String())
	require.Equal(t, "platform/v1/sifts/blocked:go-cfg/platform/BlockedSiftVersion_slice", platform.BlockedSiftsV1Key.String())
}
