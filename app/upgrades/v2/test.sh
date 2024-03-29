#!/usr/bin/env bash

#set -eo pipefail

KEYS=("dev0" "dev1" "dev2")
CHAINID="cvn_2032-2"
MONIKER="localtestnet"
# Remember to change to other types of keyring like 'file' in-case exposing to outside world,
# otherwise your balance will be wiped quickly
# The keyring test does not require private key to steal tokens from you
KEYRING="test"
# Set dedicated home directory for the cvnd instance
HOMEDIR="$HOME/.tmp2-cvnd"

# Path variables
#CONFIG=$HOMEDIR/config/config.toml
#APP_TOML=$HOMEDIR/config/app.toml
GENESIS=$HOMEDIR/config/genesis.json
TMP_GENESIS=$HOMEDIR/config/tmp_genesis.json

# validate dependencies are installed
command -v jq >/dev/null 2>&1 || {
  echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"
  exit 1
}

run_cvn_node() {
  if [ -d "$HOMEDIR" ]; then
    printf "\nAn existing folder at '%s' was found. You can choose to delete this folder and start a new local node with new keys from genesis. When declined, the existing local node is started. \n" "$HOMEDIR"
    echo "Overwrite the existing configuration and start a new local node? [y/n]"
    read -r overwrite
  else
    overwrite="Y"
  fi

  # Setup local node if overwrite is set to Yes, otherwise skip setup
  if [[ $overwrite == "y" || $overwrite == "Y" ]]; then
    # Remove the previous folder
    rm -rf "$HOMEDIR"

    # Set client config
    cvnd config keyring-backend $KEYRING --home "$HOMEDIR"
    cvnd config chain-id $CHAINID --home "$HOMEDIR"

    # If keys exist they should be deleted
    for KEY in "${KEYS[@]}"; do
      cvnd keys add "$KEY" --keyring-backend $KEYRING --home "$HOMEDIR"
    done

    # Set moniker and chain-id for cvn (Moniker can be anything, chain-id must be an integer)
    cvnd init $MONIKER -o --chain-id $CHAINID --home "$HOMEDIR"

    # Change parameter token denominations to acvnt
    jq '.app_state["staking"]["params"]["bond_denom"]="acvnt"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    jq '.app_state["crisis"]["constant_fee"]["denom"]="acvnt"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="acvnt"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    jq '.app_state["evm"]["params"]["evm_denom"]="acvnt"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    jq '.app_state["inflation"]["params"]["mint_denom"]="acvnt"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

    # Set gas limit in genesis
    jq '.consensus_params["block"]["max_gas"]="10000000"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

    # Change proposal periods to pass within a reasonable time for local testing
    sed -i.bak 's/"max_deposit_period": "172800s"/"max_deposit_period": "30s"/g' "$HOMEDIR"/config/genesis.json
    sed -i.bak 's/"voting_period": "172800s"/"voting_period": "30s"/g' "$HOMEDIR"/config/genesis.json

    # Allocate genesis accounts (cosmos formatted addresses)
    for KEY in "${KEYS[@]}"; do
      cvnd add-genesis-account "$KEY" 40000000000000000000000000acvnt --keyring-backend $KEYRING --home "$HOMEDIR"
    done
    # EIP-55 Address: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
    cvnd add-genesis-account "cvn17w0adeg64ky0daxwd2ugyuneellmjgnxp2hwdj" 40000000000000000000000000acvnt --home "$HOMEDIR"

    # bc is required to add these big numbers
    total_supply=$(echo "(${#KEYS[@]}+1) * 40000000000000000000000000" | bc)
    jq -r --arg total_supply "$total_supply" '.app_state["bank"]["supply"][0]["amount"]=$total_supply' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

    # Sign genesis transaction
    cvnd gentx "${KEYS[0]}" 30000000000000000000000000acvnt --keyring-backend $KEYRING --chain-id $CHAINID --home "$HOMEDIR"

    # Collect genesis tx
    cvnd collect-gentxs --home "$HOMEDIR"

    # Run this to ensure everything worked and that the genesis file is setup correctly
    cvnd validate-genesis --home "$HOMEDIR"
  fi

  trap 'docker stop cvn;docker rm cvn' SIGINT SIGTERM EXIT
  set -x
  docker run -it --name cvn \
    -v "$HOMEDIR/data":/root/.cvnd/data \
    -v "$HOMEDIR/config":/root/.cvnd/config \
    -p 127.0.0.1:26657:26657 -p 127.0.0.1:1317:1317 -p 127.0.0.1:8545:8545 \
    ghcr.io/cvn-network/cvn-cosmovisor:2.0.0 \
    cosmovisor run start --minimum-gas-prices=100000000acvnt \
    --rpc.laddr tcp://0.0.0.0:26657 \
    --json-rpc.address 0.0.0.0:8545 \
    --json-rpc.api eth,txpool,personal,net,debug,web3 \
    --rpc.unsafe \
    --api.enable --api.enabled-unsafe-cors
}

show_node_version() {
  curl -s http://127.0.0.1:1317/cosmos/base/tendermint/v1beta1/node_info  | jq .application_version | jq 'del(.build_deps)'
}

show_gov_module_account() {
  cvnd query auth module-account gov --output json --home "$HOMEDIR"
}

show_inflation_rate() {
  echo "inflation rate: $(cvnd query inflation inflation-rate --output json --home "$HOMEDIR")"
}

show_inflation_distribution() {
  cvnd query inflation params --output json --home "$HOMEDIR" | jq '.inflation_distribution'
}

show_base_fee() {
  cvnd query feemarket params --output json --home "$HOMEDIR" | jq -r '"base_fee: \(.base_fee)\nmin_gas_price: \(.min_gas_price)"'
}

show_slashing_signed_blocks_window() {
  cvnd query slashing params --output json --home "$HOMEDIR" | jq -r '"signed_blocks_window: \(.signed_blocks_window)"'
}

show_metadata() {
  cvnd query bank denom-metadata --output json --home "$HOMEDIR" | jq
}

deploy_soul_contract() {
  echo "deploy contract"

  if ! cd contracts; then
    echo "contracts directory not found"
    exit 1
  fi
  npm install >/dev/null 2>&1
  npx hardhat --network localhost run scripts/deploy.ts
}

submit_register_erc20_proposal_and_vote() {
  echo "submit register erc20 proposal and vote"
  local contract_address=$1
  [[ -z "$contract_address" ]] && echo "contract_address is required" && exit 1

  upgrade_height=$(cvnd status --home "$HOMEDIR" | jq -r '.SyncInfo.latest_block_height|tonumber + 20')
  echo "upgrade height = ${upgrade_height}, submitting proposal..."

  cvnd tx gov submit-legacy-proposal register-erc20 "$contract_address" \
    --title "Register ERC20 Token" \
    --deposit "10000000000000000000000acvnt" \
    --description "Register ERC20 Token" \
    --gas=auto --gas-adjustment=1.5 --gas-prices="100000000acvnt" \
    --broadcast-mode block \
    --from dev0 --home "${HOMEDIR}" --yes
  echo "create register erc20 proposal success, wait for voting..."

  sleep 5
  proposal_id=$(cvnd query gov proposals --status=voting_period --output json --home "$HOMEDIR" | jq -r '.proposals[0].id')
  echo "vote proposal id = ${proposal_id}, vote..."

  cvnd tx gov vote "${proposal_id}" yes \
    --gas=auto --gas-adjustment=1.5 --gas-prices="100000000acvnt" \
    --broadcast-mode block \
    --from dev0 --home "${HOMEDIR}" --yes
  echo "vote success, wait for proposal passed..."
}

show_erc20_token_pairs() {
  cvnd query erc20 token-pairs --output json --home "$HOMEDIR" | jq
}

submit_upgrade_proposal_and_vote() {
  echo "submit upgrade proposal and vote"

  upgrade_height=$(cvnd status --home "$HOMEDIR" | jq -r '.SyncInfo.latest_block_height|tonumber + 20')
  echo "upgrade height = ${upgrade_height}, submitting proposal..."

  cvnd tx gov submit-legacy-proposal software-upgrade "v2.0.0" \
    --title "Upgrade to v2" \
    --deposit "10000000000000000000000acvnt" \
    --description "Upgrade to v2" \
    --upgrade-height "${upgrade_height}" \
    --no-validate \
    --upgrade-info '{"binaries":{"linux/amd64":"https://github.com/cvn-network/cvn/releases/download/v2.0.0/cvnd-v2.0.0-linux-amd64"}}' \
    --gas=auto --gas-adjustment=1.5 --gas-prices="100000000acvnt" \
    --broadcast-mode block \
    --from dev0 --home "${HOMEDIR}" --yes
  echo "create upgrade proposal success, wait for voting..."

  sleep 5
  proposal_id=$(cvnd query gov proposals --status=voting_period --output json --home "$HOMEDIR" | jq -r '.proposals[0].id')
  echo "vote proposal id = ${proposal_id}, vote..."

  cvnd tx gov vote "${proposal_id}" yes \
    --gas=auto --gas-adjustment=1.5 --gas-prices="100000000acvnt" \
    --broadcast-mode block \
    --from dev0 --home "${HOMEDIR}" --yes
  echo "vote success, wait for proposal passed..."
}

show_proposal_status() {
  cvnd query gov proposals --output json --home "$HOMEDIR" | jq
}

withdraw_rewards() {
  validator_address=$(cvnd keys show dev0 --bech val --address --home "${HOMEDIR}")
  cvnd tx distribution withdraw-rewards "$validator_address" \
    --gas=auto --gas-adjustment=1.5 --gas-prices="100000000acvnt" \
    --broadcast-mode block \
    --from dev0 --home "${HOMEDIR}" --yes
}

show_balance() {
  cvnd debug addr "$(cvnd keys show dev0 --address --home "$HOMEDIR")"
  cvnd query bank balances "$(cvnd keys show dev0 --address --home "$HOMEDIR")" --output json --home "$HOMEDIR" | jq

  contract_address=$1
  from_address=$2
  [[ -z "$contract_address" ]] && exit 0
  [[ -z "$from_address" ]] && echo "from_address is required" && exit 1
  hex_balance=$(curl -s 'http://127.0.0.1:8545/' \
    -H 'Content-Type: application/json' \
    --data-raw '{"id":1,"jsonrpc":"2.0","method":"eth_call","params":[{"from":"0x0000000000000000000000000000000000000000","to":"'"$contract_address"'","data":"0x70a08231000000000000000000000000'"${from_address//0x/}"'"},"latest"]}' \
    --compressed | jq -r '.result' | sed 's/0x//' | tr '[:lower:]' '[:upper:]')
  echo "ERC20 Token Balance: $(echo "ibase=16; $hex_balance" | bc)"
}

"$@"
