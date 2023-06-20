
```shell
git clone https://github.com/cvn-network/cvn.git
cd cvn
git checkout v2.0.0

docker build -f ./cmd/cosmovisor/Dockerfile -t ghcr.io/cvn-network/cvn-cosmovisor:2.0.0 .

./app/upgrades/v2/test.sh run_cvn_node

# open another terminal && cd cvn
./app/upgrades/v2/test.sh show_inflation_rate

# open another terminal && cd cvn
./app/upgrades/v2/test.sh deploy_soul_contract
./app/upgrades/v2/test.sh submit_register_erc20_proposal_and_vote <soul-contract-address>
./app/upgrades/v2/test.sh show_proposal_status
# await for proposal to pass
./app/upgrades/v2/test.sh show_erc20_token_pairs
./app/upgrades/v2/test.sh show_metadata

# open another terminal && cd cvn
./app/upgrades/v2/test.sh submit_upgrade_proposal_and_vote
./app/upgrades/v2/test.sh show_proposal_status
# await for proposal to pass
# check node logs, await for upgrade
./app/upgrades/v2/test.sh show_inflation_rate
```

