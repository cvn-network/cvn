package keeper_test

import (
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cvn-network/cvn/v1/app"
	"github.com/cvn-network/cvn/v1/testutil"
	"github.com/cvn-network/cvn/v1/x/epochs/types"
	evm "github.com/cvn-network/cvn/v1/x/evm/types"
	"github.com/stretchr/testify/require"
)

// Test helpers
func (suite *KeeperTestSuite) DoSetupTest(t require.TestingT) {
	checkTx := false

	// init app
	suite.app = app.Setup(checkTx, nil)

	// setup context
	header := testutil.NewHeader(
		1, time.Now().UTC(), "cvn_8008-1", suite.consAddress, nil, nil,
	)
	suite.ctx = suite.app.BaseApp.NewContext(checkTx, header)

	// setup query helpers
	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.EpochsKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	// Set epoch start time and height for all epoch identifiers
	identifiers := []string{types.WeekEpochID, types.DayEpochID}
	for _, identifier := range identifiers {
		epoch, found := suite.app.EpochsKeeper.GetEpochInfo(suite.ctx, identifier)
		suite.Require().True(found)
		epoch.StartTime = suite.ctx.BlockTime()
		epoch.CurrentEpochStartHeight = suite.ctx.BlockHeight()
		suite.app.EpochsKeeper.SetEpochInfo(suite.ctx, epoch)
	}
}

func (suite *KeeperTestSuite) Commit() {
	suite.CommitAfter(time.Nanosecond)
}

func (suite *KeeperTestSuite) CommitAfter(t time.Duration) {
	var err error
	suite.ctx, err = testutil.Commit(suite.ctx, suite.app, t, nil)
	suite.Require().NoError(err)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	evm.RegisterQueryServer(queryHelper, suite.app.EvmKeeper)
	suite.queryClientEvm = evm.NewQueryClient(queryHelper)
}
