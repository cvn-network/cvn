# Upgrade v3.0.0 Test

## Environment Requirements

* OS: Mac or Linux
* Golang: 1.20+
* Docker: 24+

## Test Steps

```shell
git clone https://github.com/cvn-network/cvn.git
cd cvn
git checkout release/v2.1.x
make install
# check cvn version, should be release/v2.1.x-705ab6af
cvnd version

# build cvn-cosmovisor image
git checkout release/v3.0.x
docker build -f ./cmd/cosmovisor/Dockerfile -t ghcr.io/cvn-network/cvn-cosmovisor:3.0.0 .

# start cvn node
./app/upgrades/v3/test.sh run_cvn_node

# open another terminal && cd cvn
./app/upgrades/v3/test.sh show_node_version
./app/upgrades/v3/test.sh submit_upgrade_proposal_and_vote
# if you see `"status": "PROPOSAL_STATUS_PASSED"`, then proposal passed
./app/upgrades/v3/test.sh show_proposal_status
# await for proposal to pass
# check node logs, await for upgrade
# if you see `ERR UPGRADE "v3" NEEDED at height: xxx` and `ERR CONSENSUS FAILURE!!!`, then upgrade plan is working,
# and the node will be restarted automatically use the new binary,
# and you can see the new version in the node logs
# finally, you can see `INF executed block height=xxx`, then upgrade done
./app/upgrades/v3/test.sh show_node_version
./app/upgrades/v3/test.sh show_test_account
```

