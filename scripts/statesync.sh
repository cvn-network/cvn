#!/bin/bash
# microtick and bitcanna contributed significantly here.
# Pebbledb state sync script.
set -uxe

# Set Golang environment variables.
export GOPATH=~/go
export PATH=$PATH:~/go/bin

# Install with pebbledb 
#go mod edit -replace github.com/tendermint/tm-db=github.com/notional-labs/tm-db@136c7b6
#go mod tidy
#go install -ldflags '-w -s -X github.com/cosmos/cosmos-sdk/types.DBBackend=pebbledb' -tags pebbledb ./...

# NOTE: ABOVE YOU CAN USE ALTERNATIVE DATABASES, HERE ARE THE EXACT COMMANDS
# go install -ldflags '-w -s -X github.com/cosmos/cosmos-sdk/types.DBBackend=rocksdb' -tags rocksdb ./...
# go install -ldflags '-w -s -X github.com/cosmos/cosmos-sdk/types.DBBackend=badgerdb' -tags badgerdb ./...
# go install -ldflags '-w -s -X github.com/cosmos/cosmos-sdk/types.DBBackend=boltdb' -tags boltdb ./...

# Initialize chain.
cvnd init test --chain-id cvn_2032-1

# Get Genesis
wget https://raw.githubusercontent.com/cvn-network/cvn/main/networks/mainnet/genesis.json
mv genesis.json ~/.cvnd/config/

# Get "trust_hash" and "trust_height".
INTERVAL=1000
LATEST_HEIGHT=$(curl -s "${CVND_RPC}/block" | jq -r .result.block.header.height)
BLOCK_HEIGHT=$(("$LATEST_HEIGHT"-"$INTERVAL"))
TRUST_HASH=$(curl -s "${CVND_RPC}/block?height=$BLOCK_HEIGHT" | jq -r .result.block_id.hash)

# Print out block and transaction hash from which to sync state.
echo "trust_height: $BLOCK_HEIGHT"
echo "trust_hash: $TRUST_HASH"

# Export state sync variables.
export CVND_STATESYNC_ENABLE=true
export CVND_P2P_MAX_NUM_INBOUND_PEERS=200
export CVND_P2P_MAX_NUM_OUTBOUND_PEERS=200
#export CVND_STATESYNC_RPC_SERVERS=""
export CVND_STATESYNC_TRUST_HEIGHT=$BLOCK_HEIGHT
export CVND_STATESYNC_TRUST_HASH=$TRUST_HASH
#export CVND_P2P_SEEDS=""

# Start chain.
# Add the flag --db_backend=pebbledb if you want to use pebble.
cvnd start --x-crisis-skip-assert-invariants
