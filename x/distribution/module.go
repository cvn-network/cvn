package distribution

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/cvn-network/cvn/v2/x/distribution/keeper"
)

var _ module.AppModule = AppModule{}

// AppModule implements an application module for the distribution module.
type AppModule struct {
	distribution.AppModule
	keeper keeper.Keeper
}

func NewAppModule(appModule distribution.AppModule, keeper keeper.Keeper) AppModule {
	return AppModule{
		AppModule: appModule,
		keeper:    keeper,
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}
