package claims_test

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	"github.com/tendermint/tendermint/version"

	"github.com/cvn-network/cvn/v3/app"
	"github.com/cvn-network/cvn/v3/testutil"
	utiltx "github.com/cvn-network/cvn/v3/testutil/tx"
	"github.com/cvn-network/cvn/v3/utils"
	"github.com/cvn-network/cvn/v3/x/claims"
	"github.com/cvn-network/cvn/v3/x/claims/types"
	feemarkettypes "github.com/cvn-network/cvn/v3/x/feemarket/types"
)

type GenesisTestSuite struct {
	suite.Suite

	ctx sdk.Context

	app     *app.CVN
	genesis types.GenesisState
}

func (suite *GenesisTestSuite) SetupTest() {
	// consensus key
	consAddress := sdk.ConsAddress(utiltx.GenerateAddress().Bytes())

	suite.app = app.Setup(false, feemarkettypes.DefaultGenesisState())
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{
		Height:          1,
		ChainID:         utils.TestnetChainID + "-1",
		Time:            time.Now().UTC(),
		ProposerAddress: consAddress.Bytes(),

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

	params := types.DefaultParams()
	params.AirdropStartTime = suite.ctx.BlockTime()
	err := suite.app.ClaimsKeeper.SetParams(suite.ctx, params)
	suite.Require().NoError(err)

	stakingParams := suite.app.StakingKeeper.GetParams(suite.ctx)
	stakingParams.BondDenom = utils.BaseDenom
	suite.app.StakingKeeper.SetParams(suite.ctx, stakingParams)

	suite.genesis = *types.DefaultGenesis()
	suite.genesis.Params.AirdropStartTime = suite.ctx.BlockTime()
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

var (
	acc1 = sdk.MustAccAddressFromBech32("cvn1ydc99lv22wxzjx0r30nzcj44gvw9fmsutt8ryr")
	acc2 = sdk.MustAccAddressFromBech32("cvn1yjnzmrrr5cqr5a4xe6wscl2mvw56wj8tptws4l")
)

func (suite *GenesisTestSuite) TestClaimInitGenesis() {
	testCases := []struct {
		name     string
		genesis  types.GenesisState
		malleate func()
		expPanic bool
	}{
		{
			"default genesis",
			suite.genesis,
			func() {},
			false,
		},
		{
			"custom genesis - not all claimed",
			types.GenesisState{
				Params: suite.genesis.Params,
				ClaimsRecords: []types.ClaimsRecordAddress{
					{
						Address:                acc1.String(),
						InitialClaimableAmount: sdk.NewInt(10_000),
						ActionsCompleted:       []bool{true, false, true, true},
					},
					{
						Address:                acc2.String(),
						InitialClaimableAmount: sdk.NewInt(400),
						ActionsCompleted:       []bool{false, false, true, false},
					},
				},
			},
			func() {
				coins := sdk.NewCoins(sdk.NewCoin("acvnt", sdk.NewInt(2_800)))
				err := testutil.FundModuleAccount(suite.ctx, suite.app.BankKeeper, types.ModuleName, coins)
				suite.Require().NoError(err)
			},
			false,
		},
		{
			"custom genesis - all claimed or all unclaimed",
			types.GenesisState{
				Params: suite.genesis.Params,
				ClaimsRecords: []types.ClaimsRecordAddress{
					{
						Address:                acc1.String(),
						InitialClaimableAmount: sdk.NewInt(10_000),
						ActionsCompleted:       []bool{true, true, true, true},
					},
					{
						Address:                acc2.String(),
						InitialClaimableAmount: sdk.NewInt(400),
						ActionsCompleted:       []bool{false, false, false, false},
					},
				},
			},
			func() {
				coins := sdk.NewCoins(sdk.NewCoin("acvnt", sdk.NewInt(400)))
				err := testutil.FundModuleAccount(suite.ctx, suite.app.BankKeeper, types.ModuleName, coins)
				suite.Require().NoError(err)
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			tc.malleate()

			if tc.expPanic {
				suite.Require().Panics(func() {
					claims.InitGenesis(suite.ctx, *suite.app.ClaimsKeeper, tc.genesis)
				})
			} else {
				suite.Require().NotPanics(func() {
					claims.InitGenesis(suite.ctx, *suite.app.ClaimsKeeper, tc.genesis)
				})

				params := suite.app.ClaimsKeeper.GetParams(suite.ctx)
				suite.Require().Equal(params, tc.genesis.Params)

				claimsRecords := suite.app.ClaimsKeeper.GetClaimsRecords(suite.ctx)
				suite.Require().Equal(claimsRecords, tc.genesis.ClaimsRecords)
			}
		})
	}
}

func (suite *GenesisTestSuite) TestClaimExportGenesis() {
	suite.genesis.ClaimsRecords = []types.ClaimsRecordAddress{
		{
			Address:                acc1.String(),
			InitialClaimableAmount: sdk.NewInt(10_000),
			ActionsCompleted:       []bool{true, true, true, true},
		},
		{
			Address:                acc2.String(),
			InitialClaimableAmount: sdk.NewInt(400),
			ActionsCompleted:       []bool{false, false, false, false},
		},
	}

	coins := sdk.NewCoins(sdk.NewCoin("acvnt", sdk.NewInt(400)))
	err := testutil.FundModuleAccount(suite.ctx, suite.app.BankKeeper, types.ModuleName, coins)
	suite.Require().NoError(err)

	claims.InitGenesis(suite.ctx, *suite.app.ClaimsKeeper, suite.genesis)

	claimsRecord, found := suite.app.ClaimsKeeper.GetClaimsRecord(suite.ctx, acc2)
	suite.Require().True(found)
	suite.Require().Equal(claimsRecord, types.ClaimsRecord{
		InitialClaimableAmount: sdk.NewInt(400),
		ActionsCompleted:       []bool{false, false, false, false},
	})

	claimableAmount, remainder := suite.app.ClaimsKeeper.GetClaimableAmountForAction(suite.ctx, claimsRecord, types.ActionIBCTransfer, suite.genesis.Params)
	suite.Require().Equal(sdk.NewInt(100), claimableAmount)
	suite.Require().Equal(sdk.ZeroInt(), remainder)

	genesisExported := claims.ExportGenesis(suite.ctx, *suite.app.ClaimsKeeper)
	suite.Require().Equal(genesisExported.Params, suite.genesis.Params)
	suite.Require().Equal(genesisExported.ClaimsRecords, suite.genesis.ClaimsRecords)
}
