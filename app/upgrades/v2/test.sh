#!/usr/bin/env bash

set -eo pipefail

KEYS[0]="dev0"
KEYS[1]="dev1"
KEYS[2]="dev2"
CHAINID="cvn_2031-1"
MONIKER="localtestnet"
# Remember to change to other types of keyring like 'file' in-case exposing to outside world,
# otherwise your balance will be wiped quickly
# The keyring test does not require private key to steal tokens from you
KEYRING="test"
# Set dedicated home directory for the cvnd instance
HOMEDIR="$HOME/.tmp-cvnd"

# Path variables
CONFIG=$HOMEDIR/config/config.toml
APP_TOML=$HOMEDIR/config/app.toml
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

    # set custom pruning settings
    sed -i.bak 's/pruning = "default"/pruning = "custom"/g' "$APP_TOML"
    sed -i.bak 's/pruning-keep-recent = "0"/pruning-keep-recent = "2"/g' "$APP_TOML"
    sed -i.bak 's/pruning-interval = "0"/pruning-interval = "10"/g' "$APP_TOML"

    # Allocate genesis accounts (cosmos formatted addresses)
    for KEY in "${KEYS[@]}"; do
      cvnd add-genesis-account "$KEY" 100000000000000000000000000acvnt --keyring-backend $KEYRING --home "$HOMEDIR"
    done
    cvnd add-genesis-account "cvn17w0adeg64ky0daxwd2ugyuneellmjgnxp2hwdj" 100000000000000000000000000acvnt --home "$HOMEDIR"

    # bc is required to add these big numbers
    total_supply=$(echo "(${#KEYS[@]}+1) * 100000000000000000000000000" | bc)
    jq -r --arg total_supply "$total_supply" '.app_state["bank"]["supply"][0]["amount"]=$total_supply' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

    # Sign genesis transaction
    cvnd gentx "${KEYS[0]}" 1000000000000000000000acvnt --keyring-backend $KEYRING --chain-id $CHAINID --home "$HOMEDIR"

    # Collect genesis tx
    cvnd collect-gentxs --home "$HOMEDIR"

    # Run this to ensure everything worked and that the genesis file is setup correctly
    cvnd validate-genesis --home "$HOMEDIR"
  fi

  # Start the node (remove the --pruning=nothing flag if historical queries are not needed)
  cvnd start --minimum-gas-prices=1000000000acvnt --json-rpc.api eth,txpool,personal,net,debug,web3 --api.enable --home "$HOMEDIR"
}

show_inflation_rate() {
  echo "inflation rate: $(cvnd query inflation inflation-rate --output json --home "$HOMEDIR")"
}

show_inflation_distribution() {
  cvnd query inflation params --output json --home "$HOMEDIR" | jq '.inflation_distribution'
}

show_base_fee() {
  cvnd query feemarket base-fee --output json --home "$HOMEDIR"
}

show_metadata() {
  cvnd query bank denom-metadata --output json --home "$HOMEDIR" | jq
}

deploy_soul_contract() {
  echo "deploy contract"

  cd contracts || exit 1
  npm install
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
    --gas=auto --gas-adjustment=1.5 --gas-prices="1000000000acvnt" \
    --broadcast-mode block \
    --from dev0 --home "${HOMEDIR}" --yes
  echo "create register erc20 proposal success, wait for voting..."

  sleep 5
  proposal_id=$(cvnd query gov proposals --status=voting_period --output json --home "$HOMEDIR" | jq -r '.proposals[0].id')
  echo "vote proposal id = ${proposal_id}, vote..."

  cvnd tx gov vote "${proposal_id}" yes \
    --gas=auto --gas-adjustment=1.5 --gas-prices="1000000000acvnt" \
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

  cvnd tx gov submit-legacy-proposal software-upgrade "v2" \
    --title "Upgrade to v2" \
    --deposit "10000000000000000000000acvnt" \
    --description "Upgrade to v2" \
    --upgrade-height "${upgrade_height}" \
    --no-validate \
    --upgrade-info '{"binaries":{"linux/amd64":"https://github.com/cvn-network/cvn/releases/download/v2.0.0/cvnd-v2.0.0-linux-amd64"}}' \
    --gas=auto --gas-adjustment=1.5 --gas-prices="1000000000acvnt" \
    --broadcast-mode block \
    --from dev0 --home "${HOMEDIR}" --yes
  echo "create upgrade proposal success, wait for voting..."

  sleep 5
  proposal_id=$(cvnd query gov proposals --status=voting_period --output json --home "$HOMEDIR" | jq -r '.proposals[0].id')
  echo "vote proposal id = ${proposal_id}, vote..."

  cvnd tx gov vote "${proposal_id}" yes \
    --gas=auto --gas-adjustment=1.5 --gas-prices="1000000000acvnt" \
    --broadcast-mode block \
    --from dev0 --home "${HOMEDIR}" --yes
  echo "vote success, wait for proposal passed..."
}

show_proposal_status() {
  cvnd query gov proposals --output json --home "$HOMEDIR" | jq
}

"$@"
