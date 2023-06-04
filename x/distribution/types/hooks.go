package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var _ DistributionHooks = &MultiDistributionHooks{}

// MultiDistributionHooks defines an array of DistributionHooks whose methods are called
// in order when executing distribution functionality.
type MultiDistributionHooks []DistributionHooks

// NewMultiDistributionHooks returns a new MultiDistributionHooks instance.
func NewMultiDistributionHooks(hooks ...DistributionHooks) MultiDistributionHooks {
	return hooks
}

// AfterWithdrawDelegationRewards calls AfterWithdrawDelegationRewards method on each registered DistributionHooks.
func (h MultiDistributionHooks) AfterWithdrawDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, rewards sdk.Coins) error {
	for i := range h {
		if err := h[i].AfterWithdrawDelegationRewards(ctx, delAddr, valAddr, rewards); err != nil {
			return err
		}
	}
	return nil
}

// AfterWithdrawValidatorCommission calls AfterWithdrawValidatorCommission method on each registered DistributionHooks.
func (h MultiDistributionHooks) AfterWithdrawValidatorCommission(ctx sdk.Context, valAddr sdk.ValAddress, commission sdk.Coins) error {
	for i := range h {
		if err := h[i].AfterWithdrawValidatorCommission(ctx, valAddr, commission); err != nil {
			return err
		}
	}
	return nil
}
