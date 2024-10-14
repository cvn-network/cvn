package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"

	cvntypes "github.com/cvn-network/cvn/v3/types"
	distributiontypes "github.com/cvn-network/cvn/v3/x/distribution/types"
	erc20types "github.com/cvn-network/cvn/v3/x/erc20/types"
	"github.com/cvn-network/cvn/v3/x/gov/types"
)

var _ distributiontypes.DistributionHooks = SoulHooks{}

func NewSoulHooks(account types.AccountKeeper, bank types.BankKeeper, erc20 types.ERC20Keeper, staking types.StakingKeeper) SoulHooks {
	address := account.GetModuleAddress(govtypes.ModuleName)
	if address == nil {
		panic("gov module account has not been set")
	}
	return SoulHooks{
		govModuleAddr: address.String(),
		bank:          bank,
		erc20:         erc20,
		staking:       staking,
	}
}

type SoulHooks struct {
	govModuleAddr string
	bank          types.BankKeeper
	erc20         types.ERC20Keeper
	staking       types.StakingKeeper
}

func (m SoulHooks) AfterWithdrawDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, rewards sdk.Coins) error {
	if !delAddr.Equals(valAddr) {
		return nil
	}
	return m.mint(ctx, valAddr, sdkmath.NewIntFromBigInt(rewards.AmountOf(cvntypes.AttoCvnt).BigInt()))
}

func (m SoulHooks) AfterWithdrawValidatorCommission(ctx sdk.Context, valAddr sdk.ValAddress, commission sdk.Coins) error {
	return m.mint(ctx, valAddr, sdkmath.NewIntFromBigInt(commission.AmountOf(cvntypes.AttoCvnt).BigInt()))
}

func (m SoulHooks) mint(ctx sdk.Context, valAddr sdk.ValAddress, amount sdkmath.Int) error {
	if !amount.IsPositive() {
		return nil
	}
	coin := sdk.NewCoin(types.AttoSoult, amount)
	coins := sdk.NewCoins(coin)
	if err := m.bank.MintCoins(ctx, govtypes.ModuleName, coins); err != nil {
		return err
	}

	if err := m.bank.SendCoinsFromModuleToAccount(ctx, govtypes.ModuleName, valAddr.Bytes(), coins); err != nil {
		return err
	}

	if denomMap := m.erc20.GetDenomMap(ctx, types.AttoSoult); len(denomMap) == 0 {
		return nil
	}

	_, err := m.erc20.ConvertCoin(sdk.WrapSDKContext(ctx), &erc20types.MsgConvertCoin{
		Coin:     coin,
		Receiver: common.BytesToAddress(valAddr.Bytes()).String(),
		Sender:   sdk.AccAddress(valAddr.Bytes()).String(),
	})
	if err != nil {
		return err
	}
	return err
}
