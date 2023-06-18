package v2

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/tendermint/tendermint/libs/log"

	cvntypes "github.com/cvn-network/cvn/v2/types"
	feemarketkeeper "github.com/cvn-network/cvn/v2/x/feemarket/keeper"
	cvngovtypes "github.com/cvn-network/cvn/v2/x/gov/types"
	inflationkeeper "github.com/cvn-network/cvn/v2/x/inflation/keeper"
	inflationtypes "github.com/cvn-network/cvn/v2/x/inflation/types"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v2.0.0
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	bank bankkeeper.Keeper,
	inflation inflationkeeper.Keeper,
	slashing slashingkeeper.Keeper,
	feeMarket feemarketkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", plan.Name)
		up := NewUpgrade(logger, bank, inflation, slashing, feeMarket)
		logger.Info("running upgrade handler", "plan", plan.Name)

		up.UpdateMetadata(ctx)

		up.UpdateModuleParam(ctx)

		logger.Info("completed upgrade handler", "plan", plan.Name)
		return mm.RunMigrations(ctx, configurator, vm)
	}
}

type Upgrade struct {
	logger    log.Logger
	bank      bankkeeper.Keeper
	inflation inflationkeeper.Keeper
	slashing  slashingkeeper.Keeper
	feeMarket feemarketkeeper.Keeper
}

func NewUpgrade(logger log.Logger, bank bankkeeper.Keeper, inflation inflationkeeper.Keeper, slashing slashingkeeper.Keeper, feeMarket feemarketkeeper.Keeper) Upgrade {
	return Upgrade{
		logger:    logger,
		bank:      bank,
		inflation: inflation,
		slashing:  slashing,
		feeMarket: feeMarket,
	}
}

func (u Upgrade) UpdateMetadata(ctx sdk.Context) {
	u.logger.Info("updating acvnt denom metadata")
	u.bank.SetDenomMetaData(ctx, cvntypes.GetCvnMetadata())

	u.logger.Info("updating soul denom metadata")
	u.bank.SetDenomMetaData(ctx, cvngovtypes.GetSoulMetadata())
}

func (u Upgrade) UpdateModuleParam(ctx sdk.Context) {
	u.logger.Info("updating inflation module params")
	inflationParams := u.inflation.GetParams(ctx)
	inflationParams.InflationDistribution = inflationtypes.InflationDistribution{
		StakingRewards:  sdk.NewDecWithPrec(85, 2),
		UsageIncentives: sdk.NewDecWithPrec(5, 2),
		CommunityPool:   sdk.NewDecWithPrec(10, 2),
	}
	if err := u.inflation.SetParams(ctx, inflationParams); err != nil {
		panic(err)
	}

	u.logger.Info("updating slashing module params")
	slashingParams := u.slashing.GetParams(ctx)
	slashingParams.SignedBlocksWindow = 5000
	u.slashing.SetParams(ctx, slashingParams)

	u.logger.Info("updating fee market params")
	u.feeMarket.SetBaseFee(ctx, big.NewInt(0.1*1e9))
	return
}
