package v2

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/libs/log"

	cvntypes "github.com/cvn-network/cvn/v2/types"
	erc20keeper "github.com/cvn-network/cvn/v2/x/erc20/keeper"
	erc20types "github.com/cvn-network/cvn/v2/x/erc20/types"
	feemarketkeeper "github.com/cvn-network/cvn/v2/x/feemarket/keeper"
	cvngovtypes "github.com/cvn-network/cvn/v2/x/gov/types"
	inflationkeeper "github.com/cvn-network/cvn/v2/x/inflation/keeper"
	inflationtypes "github.com/cvn-network/cvn/v2/x/inflation/types"
)

type kvStoreKey interface {
	GetKey(storeKey string) *storetypes.KVStoreKey
}

// CreateUpgradeHandler creates an SDK upgrade handler for v2.0.0
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	bank bankkeeper.Keeper,
	inflation inflationkeeper.Keeper,
	slashing slashingkeeper.Keeper,
	feeMarket feemarketkeeper.Keeper,
	erc20 erc20keeper.Keeper,
	auth authkeeper.AccountKeeper,
	kvStoreKey kvStoreKey,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", plan.Name)
		up := NewUpgrade(logger, bank, inflation, slashing, feeMarket, erc20, auth, kvStoreKey)
		logger.Info("running upgrade handler")

		up.UpdateGovModuleAccountPermissions(ctx)

		up.UpdateSoulTokenPair(ctx)

		up.UpdateMetadata(ctx)

		up.UpdateModuleParam(ctx)

		logger.Info("completed upgrade handler")
		return mm.RunMigrations(ctx, configurator, vm)
	}
}

type Upgrade struct {
	logger     log.Logger
	auth       authkeeper.AccountKeeper
	bank       bankkeeper.Keeper
	inflation  inflationkeeper.Keeper
	slashing   slashingkeeper.Keeper
	feeMarket  feemarketkeeper.Keeper
	erc20      erc20keeper.Keeper
	kvStoreKey kvStoreKey
}

func NewUpgrade(
	logger log.Logger, bank bankkeeper.Keeper, inflation inflationkeeper.Keeper,
	slashing slashingkeeper.Keeper, feeMarket feemarketkeeper.Keeper, erc20 erc20keeper.Keeper,
	auth authkeeper.AccountKeeper, kvStoreKey kvStoreKey,
) Upgrade {
	return Upgrade{
		logger:     logger,
		bank:       bank,
		inflation:  inflation,
		slashing:   slashing,
		feeMarket:  feeMarket,
		erc20:      erc20,
		auth:       auth,
		kvStoreKey: kvStoreKey,
	}
}

func (u Upgrade) UpdateGovModuleAccountPermissions(ctx sdk.Context) {
	u.logger.Info("updating gov module account permissions")
	account := u.auth.GetModuleAccount(ctx, govtypes.ModuleName)
	govModuleAcc, ok := account.(*authtypes.ModuleAccount)
	if !ok {
		panic("gov module account is not a module account")
	}
	govModuleAcc.Permissions = append(govModuleAcc.Permissions, authtypes.Minter)
	u.auth.SetModuleAccount(ctx, govModuleAcc)
}

func (u Upgrade) UpdateSoulTokenPair(ctx sdk.Context) {
	u.logger.Info("updating soul token pair")
	u.bank.IterateAllDenomMetaData(ctx, func(metadata banktypes.Metadata) bool {
		if metadata.Symbol == "SOUL" {
			storeKey := u.kvStoreKey.GetKey(banktypes.ModuleName)
			store := ctx.KVStore(storeKey)
			store.Delete(append(banktypes.DenomMetadataPrefix, []byte(metadata.Base)...))

			tokenContract := metadata.Base[len(erc20types.ModuleName)+1:]
			tokenPairID := u.erc20.GetTokenPairID(ctx, tokenContract)
			tokenPair, found := u.erc20.GetTokenPair(ctx, tokenPairID)
			if !found {
				return true
			}
			u.erc20.DeleteTokenPair(ctx, tokenPair)

			tokenPair.Denom = erc20types.CreateBaseDenom(metadata.Symbol)
			tokenPair.ContractOwner = erc20types.OWNER_MODULE
			u.erc20.SetTokenPair(ctx, tokenPair)

			newTokenPairID := tokenPair.GetID()
			u.erc20.SetDenomMap(ctx, tokenPair.Denom, newTokenPairID)
			u.erc20.SetERC20Map(ctx, common.HexToAddress(tokenPair.Erc20Address), newTokenPairID)
			return true
		}
		return false
	})
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
	feeMarketParams := u.feeMarket.GetParams(ctx)
	feeMarketParams.BaseFee = sdk.NewInt(1e8)
	feeMarketParams.MinGasPrice = sdk.NewDecFromInt(feeMarketParams.BaseFee)
	if err := u.feeMarket.SetParams(ctx, feeMarketParams); err != nil {
		panic(err)
	}
}
