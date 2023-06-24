package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/cvn-network/cvn/v2/app"
	"github.com/cvn-network/cvn/v2/crypto/ethsecp256k1"
	"github.com/cvn-network/cvn/v2/testutil"
	"github.com/cvn-network/cvn/v2/testutil/tx"
	cvntypes "github.com/cvn-network/cvn/v2/types"
	"github.com/cvn-network/cvn/v2/x/gov/keeper"
)

func TestSoulHooksTestSuite(t *testing.T) {
	suite.Run(t, new(SoulHooksTestSuite))
}

type SoulHooksTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *app.CVN
}

func (suite *SoulHooksTestSuite) SetupTest() {
	// consensus key
	privCons, err := ethsecp256k1.GenerateKey()
	require.NoError(suite.T(), err)
	consAddress := sdk.ConsAddress(privCons.PubKey().Address())

	// init app
	suite.app = app.Setup(false, nil)
	header := testutil.NewHeader(
		1, time.Now().UTC(), "cvn_2032-1", consAddress, nil, nil,
	)
	suite.ctx = suite.app.BaseApp.NewContext(false, header)
}

func (suite *SoulHooksTestSuite) TestSoulHooks() {
	hooks := keeper.NewSoulHooks(suite.app.AccountKeeper, suite.app.BankKeeper, suite.app.Erc20Keeper, suite.app.StakingKeeper)

	accAddress, _ := tx.NewAccAddressAndKey()
	err := hooks.AfterWithdrawDelegationRewards(suite.ctx, accAddress, accAddress.Bytes(), sdk.NewCoins(sdk.NewCoin(cvntypes.AttoCvnt, sdk.NewInt(1000000000000000000))))
	require.NoError(suite.T(), err)
}
