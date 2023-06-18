package types

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

const (
	// AttoCvnt defines the default coin denomination used in CVN in:
	//
	// - Staking parameters: denomination used as stake in the dPoS chain
	// - Mint parameters: denomination minted due to fee distribution rewards
	// - Governance parameters: denomination used for spam prevention in proposal deposits
	// - Crisis parameters: constant fee denomination used for spam prevention to check broken invariant
	// - EVM parameters: denomination used for running EVM state transitions in CVN.
	AttoCvnt string = "acvnt"

	// BaseDenomUnit defines the base denomination unit for CVN.
	// 1 cvnt = 1x10^{BaseDenomUnit} acvnt
	BaseDenomUnit = 18

	// DefaultGasPrice is default gas price for evm transactions
	DefaultGasPrice = 20
)

// PowerReduction defines the default power reduction value for staking
var PowerReduction = sdkmath.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(BaseDenomUnit), nil))

// NewCvntCoin is a utility function that returns an "acvnt" coin with the given sdkmath.Int amount.
// The function will panic if the provided amount is negative.
func NewCvntCoin(amount sdkmath.Int) sdk.Coin {
	return sdk.NewCoin(AttoCvnt, amount)
}

// NewCvntDecCoin is a utility function that returns an "acvnt" decimal coin with the given sdkmath.Int amount.
// The function will panic if the provided amount is negative.
func NewCvntDecCoin(amount sdkmath.Int) sdk.DecCoin {
	return sdk.NewDecCoin(AttoCvnt, amount)
}

// NewCvntCoinInt64 is a utility function that returns an "acvnt" coin with the given int64 amount.
// The function will panic if the provided amount is negative.
func NewCvntCoinInt64(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(AttoCvnt, amount)
}

func GetCvnMetadata() banktypes.Metadata {
	return banktypes.Metadata{
		Description: "The native staking and governance token of the Conscious Network.",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    AttoCvnt,
				Exponent: 0,
			},
			{
				Denom:    "cvnt",
				Exponent: 18,
			},
		},
		Base:    AttoCvnt,
		Display: "cvnt",
		Name:    "CVN",
		Symbol:  "CVN",
	}
}
