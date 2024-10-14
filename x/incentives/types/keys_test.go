package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	utiltx "github.com/cvn-network/cvn/v3/testutil/tx"
	"github.com/cvn-network/cvn/v3/x/incentives/types"
)

func TestSplitGasMeterKey(t *testing.T) {
	contract := utiltx.GenerateAddress()
	user := utiltx.GenerateAddress()

	key := types.KeyPrefixGasMeter
	key = append(key, contract.Bytes()...)
	key = append(key, user.Bytes()...)

	contract2, user2 := types.SplitGasMeterKey(key)
	require.Equal(t, contract2, contract)
	require.Equal(t, user2, user)
}
