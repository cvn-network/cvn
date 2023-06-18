package types

import banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

const AttoSoult = "asoult"

func GetSoulMetadata() banktypes.Metadata {
	return banktypes.Metadata{
		Description: "The governance token of the Conscious Network.",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    AttoSoult,
				Exponent: 0,
			},
			{
				Denom:    "soult",
				Exponent: 18,
			},
		},
		Base:    AttoSoult,
		Display: "soult",
		Name:    "SOUL",
		Symbol:  "SOUL",
	}
}
