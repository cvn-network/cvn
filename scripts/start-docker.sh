#!/bin/bash

KEY="dev0"
CHAINID="cvn_2031-1"
MONIKER="mymoniker"
DATA_DIR=$(mktemp -d -t cvn-datadir.XXXXX)

echo "create and add new keys"
./cvnd keys add $KEY --home $DATA_DIR --no-backup --chain-id $CHAINID --algo "eth_secp256k1" --keyring-backend test
echo "init CVN with moniker=$MONIKER and chain-id=$CHAINID"
./cvnd init $MONIKER --chain-id $CHAINID --home $DATA_DIR
echo "prepare genesis: Allocate genesis accounts"
./cvnd add-genesis-account \
"$(./cvnd keys show $KEY -a --home $DATA_DIR --keyring-backend test)" 1000000000000000000acvnt \
--home $DATA_DIR --keyring-backend test
echo "prepare genesis: Sign genesis transaction"
./cvnd gentx $KEY 1000000000000000000acvnt --keyring-backend test --home $DATA_DIR --keyring-backend test --chain-id $CHAINID
echo "prepare genesis: Collect genesis tx"
./cvnd collect-gentxs --home $DATA_DIR
echo "prepare genesis: Run validate-genesis to ensure everything worked and that the genesis file is setup correctly"
./cvnd validate-genesis --home $DATA_DIR

echo "starting CVN node $i in background ..."
./cvnd start --pruning=nothing --rpc.unsafe \
--keyring-backend test --home $DATA_DIR \
>$DATA_DIR/node.log 2>&1 & disown

echo "started CVN node"
tail -f /dev/null