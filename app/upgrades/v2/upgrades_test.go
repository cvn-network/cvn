package v2_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	"github.com/tendermint/tendermint/version"

	"github.com/cvn-network/cvn/v2/app"
	v2 "github.com/cvn-network/cvn/v2/app/upgrades/v2"
	"github.com/cvn-network/cvn/v2/crypto/ethsecp256k1"
	"github.com/cvn-network/cvn/v2/types"
	erc20types "github.com/cvn-network/cvn/v2/x/erc20/types"
	feemarkettypes "github.com/cvn-network/cvn/v2/x/feemarket/types"
	cvngovtypes "github.com/cvn-network/cvn/v2/x/gov/types"
	inflationtypes "github.com/cvn-network/cvn/v2/x/inflation/types"
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
		suite.app.Erc20Keeper,
		suite.app.AccountKeeper,
		suite.app,
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

	//suite.Require().False(
	//	suite.app.BankKeeper.HasDenomMetaData(suite.ctx, "acvnt"),
	//)

	suite.Require().False(
		suite.app.BankKeeper.HasDenomMetaData(suite.ctx, "asoult"),
	)

	up.UpdateMetadata(suite.ctx)

	suite.Require().True(
		suite.app.BankKeeper.HasDenomMetaData(suite.ctx, "acvnt"),
	)
	acvntMetadataV2, found := suite.app.BankKeeper.GetDenomMetaData(suite.ctx, types.AttoCvnt)
	suite.Require().True(found)
	suite.Equal(
		types.GetCvnMetadata(),
		acvntMetadataV2,
	)

	suite.Require().True(
		suite.app.BankKeeper.HasDenomMetaData(suite.ctx, "asoult"),
	)

	soulMetadataV2, found := suite.app.BankKeeper.GetDenomMetaData(suite.ctx, cvngovtypes.AttoSoult)
	suite.Require().True(found)
	suite.Equal(
		cvngovtypes.GetSoulMetadata(),
		soulMetadataV2,
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
		suite.app.Erc20Keeper,
		suite.app.AccountKeeper,
		suite.app,
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
		sdk.NewInt(1e8).String(),
		suite.app.FeeMarketKeeper.GetParams(suite.ctx).BaseFee.String(),
	)

	suite.Require().Equal(
		sdk.NewDec(1e8).String(),
		suite.app.FeeMarketKeeper.GetParams(suite.ctx).MinGasPrice.String(),
	)
}

func (suite *UpgradeTestSuite) TestUpdateTokenPair() {
	logger := suite.ctx.Logger().With("upgrade", v2.UpgradeName)
	up := v2.NewUpgrade(
		logger,
		suite.app.BankKeeper,
		suite.app.InflationKeeper,
		suite.app.SlashingKeeper,
		suite.app.FeeMarketKeeper,
		suite.app.Erc20Keeper,
		suite.app.AccountKeeper,
		suite.app,
	)

	soulMetadataV1 := banktypes.Metadata{
		Description: "Cosmos coin token representation of 0x5FbDB2315678afecb367f032d93F642f64180aa3",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "erc20/0x5FbDB2315678afecb367f032d93F642f64180aa3",
				Exponent: 0,
			},
			{
				Denom:    "SOUL",
				Exponent: 18,
			},
		},
		Base:    "erc20/0x5FbDB2315678afecb367f032d93F642f64180aa3",
		Display: "SOUL",
		Name:    "erc20/0x5FbDB2315678afecb367f032d93F642f64180aa3",
		Symbol:  "SOUL",
	}
	suite.app.BankKeeper.SetDenomMetaData(suite.ctx, soulMetadataV1)

	contractAddr := common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3")
	pair := erc20types.NewTokenPair(contractAddr, soulMetadataV1.Base, erc20types.OWNER_EXTERNAL)
	suite.app.Erc20Keeper.SetTokenPair(suite.ctx, pair)
	pairID := pair.GetID()
	suite.app.Erc20Keeper.SetDenomMap(suite.ctx, pair.Denom, pairID)
	suite.app.Erc20Keeper.SetERC20Map(suite.ctx, contractAddr, pairID)

	up.UpdateSoulTokenPair(suite.ctx)

	hasDenomMetaData := suite.app.BankKeeper.HasDenomMetaData(suite.ctx, soulMetadataV1.Base)
	suite.Require().False(hasDenomMetaData)

	pairs := suite.app.Erc20Keeper.GetTokenPairs(suite.ctx)
	suite.Require().Equal(1, len(pairs))
	pair = pairs[0]
	suite.Require().Equal(erc20types.OWNER_MODULE, pair.ContractOwner)
	suite.Require().Equal(erc20types.CreateBaseDenom("SOUL"), pair.Denom)
	suite.Require().Equal(contractAddr.String(), pair.Erc20Address)
	suite.Require().Equal(true, pair.Enabled)
	suite.Require().NotEqual(pairID, pair.GetID())

	suite.Equal(
		pair.GetID(),
		suite.app.Erc20Keeper.GetERC20Map(suite.ctx, contractAddr),
	)
	suite.Equal(
		pair.GetID(),
		suite.app.Erc20Keeper.GetDenomMap(suite.ctx, pair.Denom),
	)
	suite.Equal(
		pair.GetID(),
		suite.app.Erc20Keeper.GetTokenPairID(suite.ctx, contractAddr.String()),
	)
	tokenPair, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, pair.GetID())
	suite.Require().True(found)
	suite.Require().Equal(pair, tokenPair)
}

func TestUpgradeTestSuite(t *testing.T) {
	s := new(UpgradeTestSuite)
	suite.Run(t, s)
}
