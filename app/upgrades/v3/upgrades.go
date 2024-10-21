package v3

import (
	_ "embed"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/libs/log"
)

//go:embed cvn.json
var MigrateJSON []byte

// CreateUpgradeHandler creates an SDK upgrade handler for v2.0.0
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	bank bankkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", plan.Name)
		up := NewUpgrade(logger, bank)
		logger.Info("running upgrade handler")

		if err := up.MigrateCVNToken(ctx); err != nil {
			logger.Error("failed to migrate cvn token", "error", err)
		}

		logger.Info("completed upgrade handler")
		return mm.RunMigrations(ctx, configurator, vm)
	}
}

type Upgrade struct {
	logger log.Logger
	bank   bankkeeper.Keeper
}

func NewUpgrade(logger log.Logger, bank bankkeeper.Keeper) Upgrade {
	return Upgrade{
		logger: logger,
		bank:   bank,
	}
}

func (u Upgrade) MigrateCVNToken(ctx sdk.Context) error {
	migrates, err := ReadMigrates()
	if err != nil {
		return err
	}
	for _, migrate := range migrates {
		holderAcc := sdk.AccAddress(migrate.Holder.Bytes())
		toAcc := sdk.AccAddress(migrate.To.Bytes())
		balances := u.bank.GetAllBalances(ctx, holderAcc)
		if balances.IsZero() {
			continue
		}
		if err = u.bank.SendCoins(ctx, holderAcc, toAcc, balances); err != nil {
			return err
		}
	}
	return nil
}

type Migrate struct {
	Holder common.Address `json:"holder"`
	To     common.Address `json:"to"`
}

func ReadMigrates() ([]Migrate, error) {
	var migrates []Migrate
	if err := json.Unmarshal(MigrateJSON, &migrates); err != nil {
		return nil, err
	}
	return migrates, nil
}
