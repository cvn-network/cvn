package types

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type InflationTestSuite struct {
	suite.Suite
}

func TestInflationSuite(t *testing.T) {
	suite.Run(t, new(InflationTestSuite))
}

func (suite *InflationTestSuite) TestCalculateEpochMintProvision() {
	bondingParams := DefaultParams()
	//bondingParams.ExponentialCalculation.MaxVariance = sdk.NewDecWithPrec(40, 2)
	epochsPerPeriod := int64(365)

	testCases := []struct {
		name              string
		params            Params
		period            uint64
		bondedRatio       sdk.Dec
		expEpochProvision sdk.Dec
		expPass           bool
	}{
		{
			"pass - initial perid",
			DefaultParams(),
			uint64(0),
			sdk.OneDec(),
			sdk.MustNewDecFromStr("123287" + "671232876712328767.000000000000000000"),
			true,
		},
		{
			"pass - period 1",
			DefaultParams(),
			uint64(1),
			sdk.OneDec(),
			sdk.MustNewDecFromStr("68493150684931506849315.000000000000000000"),
			true,
		},
		{
			"pass - period 2",
			DefaultParams(),
			uint64(2),
			sdk.OneDec(),
			sdk.MustNewDecFromStr("41095890410958904109589.000000000000000000"),
			true,
		},
		{
			"pass - period 3",
			DefaultParams(),
			uint64(3),
			sdk.OneDec(),
			sdk.MustNewDecFromStr("27397260273972602739726.000000000000000000"),
			true,
		},
		{
			"pass - period 20",
			DefaultParams(),
			uint64(20),
			sdk.OneDec(),
			sdk.MustNewDecFromStr("13698734649240154082192.000000000000000000"),
			true,
		},
		{
			"pass - period 21",
			DefaultParams(),
			uint64(21),
			sdk.OneDec(),
			sdk.MustNewDecFromStr("13698682393113227726027.000000000000000000"),
			true,
		},
		{
			"pass - period 60",
			DefaultParams(),
			uint64(60),
			sdk.OneDec(),
			sdk.MustNewDecFromStr("13698" + "630136986301479452.000000000000000000"),
			true,
		},
		{
			"pass - 0 percent bonding - initial period",
			bondingParams,
			uint64(0),
			sdk.ZeroDec(),
			sdk.MustNewDecFromStr("123287" + "671232876712328767.000000000000000000"),
			true,
		},
		{
			"pass - 0 percent bonding - period 1",
			bondingParams,
			uint64(1),
			sdk.ZeroDec(),
			sdk.MustNewDecFromStr("68493150684931506849315.000000000000000000"),
			true,
		},
		{
			"pass - 0 percent bonding - period 2",
			bondingParams,
			uint64(2),
			sdk.ZeroDec(),
			sdk.MustNewDecFromStr("41095890410958904109589.000000000000000000"),
			true,
		},
		{
			"pass - 0 percent bonding - period 3",
			bondingParams,
			uint64(3),
			sdk.ZeroDec(),
			sdk.MustNewDecFromStr("27397260273972602739726.000000000000000000"),
			true,
		},
		{
			"pass - 0 percent bonding - period 20",
			bondingParams,
			uint64(20),
			sdk.ZeroDec(),
			sdk.MustNewDecFromStr("13698734649240154082192.000000000000000000"),
			true,
		},
		{
			"pass - 0 percent bonding - period 21",
			bondingParams,
			uint64(21),
			sdk.ZeroDec(),
			sdk.MustNewDecFromStr("13698682393113227726027.000000000000000000"),
			true,
		},
		{
			"pass - 0 percent bonding - period 60",
			bondingParams,
			uint64(60),
			sdk.ZeroDec(),
			sdk.MustNewDecFromStr("13698" + "630136986301479452.000000000000000000"),
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			epochMintProvisions := CalculateEpochMintProvision(
				tc.params,
				tc.period,
				epochsPerPeriod,
				tc.bondedRatio,
			)

			suite.Require().Equal(tc.expEpochProvision, epochMintProvisions)
		})
	}
}
