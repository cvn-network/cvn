# Upgrade v2.0.0 Test

## Environment Requirements

* OS: Mac or Linux
* Golang: 1.20+
* Node.js: 14+
* Docker: 24+

## Test Steps

```shell
git clone https://github.com/cvn-network/cvn.git
cd cvn
git checkout release/v1.0.x
make install
# check cvn version, should be release/v1.0.x-101a8cfb
cvnd version

# build cvn-cosmovisor image
git checkout release/v2.0.x
docker build -f ./cmd/cosmovisor/Dockerfile -t ghcr.io/cvn-network/cvn-cosmovisor:2.0.0 .

./app/upgrades/v2/test.sh run_cvn_node
./app/upgrades/v2/test.sh show_node_version

# open another terminal && cd cvn
./app/upgrades/v2/test.sh show_inflation_rate
./app/upgrades/v2/test.sh show_inflation_distribution
./app/upgrades/v2/test.sh show_base_fee
./app/upgrades/v2/test.sh show_slashing_signed_blocks_window
./app/upgrades/v2/test.sh show_gov_module_account

# open another terminal && cd cvn
./app/upgrades/v2/test.sh deploy_soul_contract
./app/upgrades/v2/test.sh submit_register_erc20_proposal_and_vote <soul-contract-address>
# if you see `"status": "PROPOSAL_STATUS_PASSED"`, then proposal passed
./app/upgrades/v2/test.sh show_proposal_status
# await for proposal to pass
./app/upgrades/v2/test.sh show_erc20_token_pairs
./app/upgrades/v2/test.sh show_metadata

# open another terminal && cd cvn
./app/upgrades/v2/test.sh submit_upgrade_proposal_and_vote
# if you see `"status": "PROPOSAL_STATUS_PASSED"`, then proposal passed
./app/upgrades/v2/test.sh show_proposal_status
# await for proposal to pass
# check node logs, await for upgrade
# if you see `ERR UPGRADE "v2" NEEDED at height: xxx` and `ERR CONSENSUS FAILURE!!!`, then upgrade plan is working,
# and the node will be restarted automatically use the new binary,
# and you can see the new version in the node logs
# finally, you can see `INF executed block height=xxx`, then upgrade done

# open another terminal && cd cvn
make install
# check cvn version, should be release/v2.0.x-xxx
cvnd version
./app/upgrades/v2/test.sh show_inflation_rate
./app/upgrades/v2/test.sh show_inflation_distribution
./app/upgrades/v2/test.sh show_base_fee
./app/upgrades/v2/test.sh show_slashing_signed_blocks_window

./app/upgrades/v2/test.sh show_erc20_token_pairs
./app/upgrades/v2/test.sh show_metadata

# open another terminal && cd cvn
./app/upgrades/v2/test.sh withdraw_rewards
./app/upgrades/v2/test.sh show_balance <soul-contract-address> <from-address>
```

