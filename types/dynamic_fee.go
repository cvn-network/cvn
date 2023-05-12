package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	evmostypes "github.com/evmos/evmos/v12/types"
)

// HasDynamicFeeExtensionOption returns true if the tx implements the `ExtensionOptionDynamicFeeTx` extension option.
func HasDynamicFeeExtensionOption(any *codectypes.Any) bool {
	_, ok := any.GetCachedValue().(*evmostypes.ExtensionOptionDynamicFeeTx)
	return ok
}
