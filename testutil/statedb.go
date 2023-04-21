package testutil

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cvn-network/cvn/v1/app/ante/evm"
	"github.com/cvn-network/cvn/v1/x/evm/statedb"
	"github.com/ethereum/go-ethereum/common"
)

// NewStateDB returns a new StateDB for testing purposes.
func NewStateDB(ctx sdk.Context, evmKeeper evm.EVMKeeper) *statedb.StateDB {
	return statedb.New(ctx, evmKeeper, statedb.NewEmptyTxConfig(common.BytesToHash(ctx.HeaderHash().Bytes())))
}
