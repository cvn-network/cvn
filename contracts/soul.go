package contracts

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"

	evmtypes "github.com/cvn-network/cvn/v2/x/evm/types"
)

var (
	//go:embed compiled_contracts/SOUL.json
	SoulJSON []byte //nolint: golint

	// SoulContract is the compiled Soul contract
	SoulContract evmtypes.CompiledContract

	// SoulAddress is the Soul address
	SoulAddress common.Address
)

func init() {
	SoulAddress = common.HexToAddress("")

	err := json.Unmarshal(SoulJSON, &SoulContract)
	if err != nil {
		panic(err)
	}

	if len(SoulContract.Bin) == 0 {
		panic("load contract failed")
	}
}
