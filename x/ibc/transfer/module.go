package transfer

import (
	"github.com/cosmos/cosmos-sdk/types/module"

	ibctransfer "github.com/cosmos/ibc-go/v6/modules/apps/transfer"
	"github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	"github.com/cvn-network/cvn/v1/x/ibc/transfer/keeper"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic embeds the IBC Transfer AppModuleBasic
type AppModuleBasic struct {
	*ibctransfer.AppModuleBasic
}

// AppModule represents the AppModule for this module
type AppModule struct {
	*ibctransfer.AppModule
	keeper keeper.Keeper
}

// NewAppModule creates a new 20-transfer module
func NewAppModule(k keeper.Keeper) AppModule {
	am := ibctransfer.NewAppModule(*k.Keeper)
	return AppModule{
		AppModule: &am,
		keeper:    k,
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	// Override Transfer Msg Server
	types.RegisterMsgServer(cfg.MsgServer(), am.keeper)
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}
