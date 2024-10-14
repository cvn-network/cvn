package v3

import (
	_ "embed"
	"encoding/json"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/cvn-network/cvn/v3/types"
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
		up := NewUpgrade(logger, bank, "")
		logger.Info("running upgrade handler")

		if err := up.MigrateCVNToken(ctx); err != nil {
			logger.Error("failed to migrate cvn token", "error", err)
		}

		logger.Info("completed upgrade handler")
		return mm.RunMigrations(ctx, configurator, vm)
	}
}

type Upgrade struct {
	logger    log.Logger
	bank      bankkeeper.Keeper
	toAddress string
}

func NewUpgrade(logger log.Logger, bank bankkeeper.Keeper, toAddr string) Upgrade {
	return Upgrade{
		logger:    logger,
		bank:      bank,
		toAddress: toAddr,
	}
}

func (u Upgrade) MigrateCVNToken(ctx sdk.Context) error {
	migrates, err := ReadMigrates()
	if err != nil {
		return err
	}
	hexToAddress := common.HexToAddress(u.toAddress)
	for _, migrate := range migrates {
		balances := u.bank.GetAllBalances(ctx, migrate.Holder)
		if balances.IsZero() {
			continue
		}
		cvnAmount := balances.AmountOf(types.AttoCvnt)
		if !cvnAmount.Equal(migrate.Value) {
			u.logger.Info("expected balance does not match", "holder", migrate.Holder.String(), "expected", migrate.Value.String(), "actual", cvnAmount.String())
		}
		if err = u.bank.SendCoins(ctx, migrate.Holder, hexToAddress.Bytes(), balances); err != nil {
			return err
		}
	}
	return nil
}

type Migrate struct {
	Holder sdk.AccAddress `json:"holder"`
	Value  sdkmath.Int    `json:"value"`
}

func ReadMigrates() ([]Migrate, error) {
	var ms []struct {
		Holder string  `json:"holder"`
		Value  float64 `json:"value"`
	}
	if err := json.Unmarshal(MigrateJSON, &ms); err != nil {
		return nil, err
	}
	migrates := make([]Migrate, 0, len(ms))
	for _, m := range ms {
		value, err := sdkmath.LegacyNewDecFromStr(fmt.Sprintf("%f", m.Value))
		if err != nil {
			return nil, err
		}
		migrates = append(migrates, Migrate{
			Holder: common.HexToAddress(m.Holder).Bytes(),
			Value:  value.MulInt64(1e18).RoundInt(),
		})
	}
	return migrates, nil
}
