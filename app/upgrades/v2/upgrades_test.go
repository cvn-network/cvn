package v2_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	"github.com/tendermint/tendermint/version"

	"github.com/cvn-network/cvn/v1/app"
	v2 "github.com/cvn-network/cvn/v1/app/upgrades/v2"
	"github.com/cvn-network/cvn/v1/crypto/ethsecp256k1"
	feemarkettypes "github.com/cvn-network/cvn/v1/x/feemarket/types"
	inflationtypes "github.com/cvn-network/cvn/v1/x/inflation/types"
)

type UpgradeTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	app         *app.CVN
	consAddress sdk.ConsAddress
}

func (suite *UpgradeTestSuite) SetupTest() {
	// consensus key
	consensusPrivKey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	suite.consAddress = sdk.ConsAddress(consensusPrivKey.PubKey().Address())

	// NOTE: this is the new binary, not the old one.
	suite.app = app.Setup(false, feemarkettypes.DefaultGenesisState())
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{
		Height:          1,
		ChainID:         "test-chain-1",
		Time:            time.Now(),
		ProposerAddress: suite.consAddress.Bytes(),

		Version: tmversion.Consensus{
			Block: version.BlockProtocol,
		},
		LastBlockId: tmproto.BlockID{
			Hash: tmhash.Sum([]byte("block_id")),
			PartSetHeader: tmproto.PartSetHeader{
				Total: 11,
				Hash:  tmhash.Sum([]byte("partset_header")),
			},
		},
		AppHash:            tmhash.Sum([]byte("app")),
		DataHash:           tmhash.Sum([]byte("data")),
		EvidenceHash:       tmhash.Sum([]byte("evidence")),
		ValidatorsHash:     tmhash.Sum([]byte("validators")),
		NextValidatorsHash: tmhash.Sum([]byte("next_validators")),
		ConsensusHash:      tmhash.Sum([]byte("consensus")),
		LastResultsHash:    tmhash.Sum([]byte("last_result")),
	})

	cp := suite.app.BaseApp.GetConsensusParams(suite.ctx)
	suite.ctx = suite.ctx.WithConsensusParams(cp)
}

func (suite *UpgradeTestSuite) TestUpdateMetadata() {
	logger := suite.ctx.Logger().With("upgrade", v2.UpgradeName)
	up := v2.NewUpgrade(
		logger,
		suite.app.BankKeeper,
		suite.app.InflationKeeper,
		suite.app.SlashingKeeper,
		suite.app.FeeMarketKeeper,
	)

	suite.Require().Equal(
		inflationtypes.DefaultInflationDistribution.StakingRewards.String(),
		suite.app.InflationKeeper.GetParams(suite.ctx).InflationDistribution.StakingRewards.String(),
	)

	suite.Require().Equal(
		inflationtypes.DefaultInflationDistribution.UsageIncentives.String(),
		suite.app.InflationKeeper.GetParams(suite.ctx).InflationDistribution.UsageIncentives.String(),
	)

	suite.Require().Equal(
		inflationtypes.DefaultInflationDistribution.CommunityPool.String(),
		suite.app.InflationKeeper.GetParams(suite.ctx).InflationDistribution.CommunityPool.String(),
	)

	suite.Require().False(
		suite.app.BankKeeper.HasDenomMetaData(suite.ctx, "acvnt"),
	)

	suite.Require().False(
		suite.app.BankKeeper.HasDenomMetaData(suite.ctx, "asoult"),
	)

	up.UpdateMetadata(suite.ctx)

	suite.Require().True(
		suite.app.BankKeeper.HasDenomMetaData(suite.ctx, "acvnt"),
	)

	suite.Require().True(
		suite.app.BankKeeper.HasDenomMetaData(suite.ctx, "asoult"),
	)
}

func (suite *UpgradeTestSuite) TestUpdateParams() {
	logger := suite.ctx.Logger().With("upgrade", v2.UpgradeName)
	up := v2.NewUpgrade(
		logger,
		suite.app.BankKeeper,
		suite.app.InflationKeeper,
		suite.app.SlashingKeeper,
		suite.app.FeeMarketKeeper,
	)

	suite.Require().Equal(
		int64(100),
		suite.app.SlashingKeeper.GetParams(suite.ctx).SignedBlocksWindow,
	)

	suite.Require().Equal(
		feemarkettypes.DefaultParams().BaseFee,
		suite.app.FeeMarketKeeper.GetParams(suite.ctx).BaseFee,
	)

	up.UpdateModuleParam(suite.ctx)

	suite.Require().Equal(
		"0.850000000000000000",
		suite.app.InflationKeeper.GetParams(suite.ctx).InflationDistribution.StakingRewards.String(),
	)

	suite.Require().Equal(
		"0.050000000000000000",
		suite.app.InflationKeeper.GetParams(suite.ctx).InflationDistribution.UsageIncentives.String(),
	)

	suite.Require().Equal(
		"0.100000000000000000",
		suite.app.InflationKeeper.GetParams(suite.ctx).InflationDistribution.CommunityPool.String(),
	)

	suite.Require().Equal(
		int64(5000),
		suite.app.SlashingKeeper.GetParams(suite.ctx).SignedBlocksWindow,
	)

	suite.Require().Equal(
		sdk.NewInt(1e8),
		suite.app.FeeMarketKeeper.GetParams(suite.ctx).BaseFee,
	)
}

func TestUpgradeTestSuite(t *testing.T) {
	s := new(UpgradeTestSuite)
	suite.Run(t, s)
}
