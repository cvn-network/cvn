package v3_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cvn-network/cvn/v3/types"
	inflationtypes "github.com/cvn-network/cvn/v3/x/inflation/types"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	"github.com/tendermint/tendermint/version"

	"github.com/cvn-network/cvn/v3/app"
	v3 "github.com/cvn-network/cvn/v3/app/upgrades/v3"
	"github.com/cvn-network/cvn/v3/crypto/ethsecp256k1"
	feemarkettypes "github.com/cvn-network/cvn/v3/x/feemarket/types"
)

type UpgradeTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	app         *app.CVN
	consAddress sdk.ConsAddress
}

func TestUpgradeTestSuite(t *testing.T) {
	s := new(UpgradeTestSuite)
	suite.Run(t, s)
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
		ChainID:         "cvn_2032-1",
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

func (suite *UpgradeTestSuite) TestUpgrade() {
	migrates, err := v3.ReadMigrates()
	suite.NoError(err)
	for _, migrate := range migrates {
		coins := sdk.NewCoins(sdk.NewCoin(types.AttoCvnt, sdk.NewIntFromUint64(tmrand.Uint64())))
		suite.NoError(suite.app.BankKeeper.MintCoins(suite.ctx, inflationtypes.ModuleName, coins))
		suite.NoError(suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, inflationtypes.ModuleName, migrate.Holder.Bytes(), coins))
	}

	logger := suite.ctx.Logger().With("upgrade", v3.UpgradeName)
	upgrade := v3.NewUpgrade(
		logger,
		suite.app.BankKeeper,
	)

	suite.NoError(upgrade.MigrateCVNToken(suite.ctx))
}
