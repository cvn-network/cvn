#!/usr/bin/env bash

set -eo pipefail

solidity_contracts=("SOUL.sol")
for contract in "${solidity_contracts[@]}"; do
  echo "===> Compiling contract: $contract"
  [[ -d ./out ]] && rm -r ./out
  [[ ! -d ./out ]] && mkdir -p ./out
  npx solc --bin --abi -p -o ./out --include-path node_modules/ --base-path . "$contract"
  contract_name=${contract%%.sol}
  cat >"./compiled_contracts/${contract_name}.json" <<EOF
{
  "abi": $(jq 'tojson' "./out/${contract_name}_sol_${contract_name}.abi"),
  "bin": "$(cat "./out/${contract_name}_sol_${contract_name}.bin")",
  "contractName": "${contract_name}"
}
EOF
done

[[ -d ./out ]] && rm -r ./out
