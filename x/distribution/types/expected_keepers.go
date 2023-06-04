package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DistributionHooks interface {
	AfterWithdrawDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, rewards sdk.Coins) error
	AfterWithdrawValidatorCommission(ctx sdk.Context, valAddr sdk.ValAddress, commission sdk.Coins) error
}
