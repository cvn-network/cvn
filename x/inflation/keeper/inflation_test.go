package keeper_test

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	cvntypes "github.com/cvn-network/cvn/v3/types"
	incentivestypes "github.com/cvn-network/cvn/v3/x/incentives/types"
	"github.com/cvn-network/cvn/v3/x/inflation/types"
)

func (suite *KeeperTestSuite) TestMintAndAllocateInflation() {
	testCases := []struct {
		name                  string
		mintCoin              sdk.Coin
		malleate              func()
		expStakingRewardAmt   sdk.Coin
		expUsageIncentivesAmt sdk.Coin
		expCommunityPoolAmt   sdk.DecCoins
		expPass               bool
	}{
		{
			"pass",
			sdk.NewCoin(denomMint, sdk.NewInt(1_000_000)),
			func() {},
			sdk.NewCoin(denomMint, sdk.NewInt(850000)),
			sdk.NewCoin(denomMint, sdk.NewInt(50000)),
			sdk.NewDecCoins(sdk.NewDecCoin(denomMint, sdk.NewInt(100000))),
			true,
		},
		{
			"pass - no coins minted ",
			sdk.NewCoin(denomMint, sdk.ZeroInt()),
			func() {},
			sdk.NewCoin(denomMint, sdk.ZeroInt()),
			sdk.NewCoin(denomMint, sdk.ZeroInt()),
			sdk.DecCoins(nil),
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			tc.malleate()

			_, _, _, err := suite.app.InflationKeeper.MintAndAllocateInflation(suite.ctx, tc.mintCoin, types.DefaultParams())

			// Get balances
			balanceModule := suite.app.BankKeeper.GetBalance(
				suite.ctx,
				suite.app.AccountKeeper.GetModuleAddress(types.ModuleName),
				denomMint,
			)

			feeCollector := suite.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
			balanceStakingRewards := suite.app.BankKeeper.GetBalance(
				suite.ctx,
				feeCollector,
				denomMint,
			)

			incentives := suite.app.AccountKeeper.GetModuleAddress(incentivestypes.ModuleName)
			balanceUsageIncentives := suite.app.BankKeeper.GetBalance(
				suite.ctx,
				incentives,
				denomMint,
			)

			balanceCommunityPool := suite.app.DistrKeeper.GetFeePoolCommunityCoins(suite.ctx)

			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().True(balanceModule.IsZero())
				suite.Require().Equal(tc.expStakingRewardAmt, balanceStakingRewards)
				suite.Require().Equal(tc.expUsageIncentivesAmt, balanceUsageIncentives)
				suite.Require().Equal(tc.expCommunityPoolAmt, balanceCommunityPool)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGetCirculatingSupplyAndInflationRate() {
	// the total bonded tokens for the 2 accounts initialized on the setup
	bondedAmt := sdkmath.NewInt(1000100000000000000)

	testCases := []struct {
		name             string
		bankSupply       sdkmath.Int
		malleate         func()
		expInflationRate sdk.Dec
	}{
		{
			"no epochs per period",
			sdk.TokensFromConsensusPower(200_000_000, cvntypes.PowerReduction).Sub(bondedAmt),
			func() {
				suite.app.InflationKeeper.SetEpochsPerPeriod(suite.ctx, 0)
			},
			sdk.ZeroDec(),
		},
		{
			"high supply",
			sdk.TokensFromConsensusPower(90_000_000, cvntypes.PowerReduction).Sub(bondedAmt),
			func() {},
			sdk.MustNewDecFromStr("50.000000000000000000"),
		},
		{
			"low supply",
			sdk.TokensFromConsensusPower(30_000_000, cvntypes.PowerReduction).Sub(bondedAmt),
			func() {},
			sdk.MustNewDecFromStr("150.000000000000000000"),
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			suite.ctx = suite.ctx.WithChainID("cvn_2032-1")
			tc.malleate()

			// Mint coins to increase supply
			mintCoin := sdk.NewCoin(
				types.DefaultInflationDenom,
				tc.bankSupply,
			)
			err := suite.app.InflationKeeper.MintCoins(suite.ctx, mintCoin)
			suite.Require().NoError(err)

			circulatingSupply := s.app.InflationKeeper.GetCirculatingSupply(suite.ctx, types.DefaultInflationDenom)
			suite.Require().Equal(sdk.NewDecCoinFromCoin(mintCoin.AddAmount(bondedAmt)).Amount, circulatingSupply)

			inflationRate := s.app.InflationKeeper.GetInflationRate(suite.ctx, types.DefaultInflationDenom)
			suite.Require().Equal(tc.expInflationRate, inflationRate)
		})
	}
}

func (suite *KeeperTestSuite) TestBondedRatio() {
	testCases := []struct {
		name         string
		malleate     func()
		expBondRatio sdk.Dec
	}{
		{
			"bonded ratio",
			func() {},
			sdk.MustNewDecFromStr("0.999900009999000099"),
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			tc.malleate()

			bondRatio := suite.app.InflationKeeper.BondedRatio(suite.ctx)
			suite.Require().Equal(tc.expBondRatio, bondRatio)
		})
	}
}
