package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	distributiontypes "github.com/cvn-network/cvn/v3/x/distribution/types"
)

type Hooks struct {
	keeper.Hooks
	k Keeper
}

var _ stakingtypes.StakingHooks = Hooks{}

// Hooks create new distribution hooks
func (k Keeper) Hooks() Hooks {
	return Hooks{
		Hooks: k.Keeper.Hooks(),
		k:     k,
	}
}

// AfterValidatorRemoved performs clean up after a validator is removed
func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) error {
	// fetch outstanding
	outstanding := h.k.GetValidatorOutstandingRewardsCoins(ctx, valAddr)

	// force-withdraw commission
	commission := h.k.GetValidatorAccumulatedCommission(ctx, valAddr).Commission
	if !commission.IsZero() {
		// subtract from outstanding
		outstanding = outstanding.Sub(commission)

		// split into integral & remainder
		coins, remainder := commission.TruncateDecimal()

		// remainder to community pool
		feePool := h.k.GetFeePool(ctx)
		feePool.CommunityPool = feePool.CommunityPool.Add(remainder...)
		h.k.SetFeePool(ctx, feePool)

		// add to validator account
		if !coins.IsZero() {
			accAddr := sdk.AccAddress(valAddr)
			withdrawAddr := h.k.GetDelegatorWithdrawAddr(ctx, accAddr)

			if err := h.k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawAddr, coins); err != nil {
				return err
			}

			if err := h.k.AfterWithdrawValidatorCommission(ctx, valAddr, coins); err != nil {
				return err
			}
		}
	}

	// Add outstanding to community pool
	// The validator is removed only after it has no more delegations.
	// This operation sends only the remaining dust to the community pool.
	feePool := h.k.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(outstanding...)
	h.k.SetFeePool(ctx, feePool)

	// delete outstanding
	h.k.DeleteValidatorOutstandingRewards(ctx, valAddr)

	// remove commission record
	h.k.DeleteValidatorAccumulatedCommission(ctx, valAddr)

	// clear slashes
	h.k.DeleteValidatorSlashEvents(ctx, valAddr)

	// clear historical rewards
	h.k.DeleteValidatorHistoricalRewards(ctx, valAddr)

	// clear current rewards
	h.k.DeleteValidatorCurrentRewards(ctx, valAddr)

	return nil
}

// BeforeDelegationSharesModified withdraw delegation rewards (which also increments period)
func (h Hooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	val := h.k.stakingKeeper.Validator(ctx, valAddr)
	del := h.k.stakingKeeper.Delegation(ctx, delAddr, valAddr)

	rewards, err := h.k.withdrawDelegationRewards(ctx, val, del)
	if err != nil {
		return err
	}

	if err = h.k.AfterWithdrawDelegationRewards(ctx, delAddr, valAddr, rewards); err != nil {
		return err
	}

	return nil
}

var _ distributiontypes.DistributionHooks = Keeper{}

func (k Keeper) AfterWithdrawDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, rewards sdk.Coins) error {
	if k.hooks != nil {
		return k.hooks.AfterWithdrawDelegationRewards(ctx, delAddr, valAddr, rewards)
	}
	return nil
}

func (k Keeper) AfterWithdrawValidatorCommission(ctx sdk.Context, valAddr sdk.ValAddress, commission sdk.Coins) error {
	if k.hooks != nil {
		return k.hooks.AfterWithdrawValidatorCommission(ctx, valAddr, commission)
	}
	return nil
}
