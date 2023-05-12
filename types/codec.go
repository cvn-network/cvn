package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	evmostypes "github.com/evmos/evmos/v12/types"
)

// RegisterInterfaces registers the tendermint concrete client-related
// implementations and interfaces.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*authtypes.AccountI)(nil),
		&evmostypes.EthAccount{},
	)
	registry.RegisterImplementations(
		(*authtypes.GenesisAccount)(nil),
		&evmostypes.EthAccount{},
	)
	registry.RegisterImplementations(
		(*tx.TxExtensionOptionI)(nil),
		&evmostypes.ExtensionOptionsWeb3Tx{},
		&evmostypes.ExtensionOptionDynamicFeeTx{},
	)
}
