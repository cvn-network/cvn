# Build image for a local testnet
localnet-build:
	docker build --no-cache --tag cvn/node -f ./networks/local/Dockerfile .

# Start a 4-node testnet locally
localnet-start: localnet-stop localnet-build
	@if ! [ -f build/node0/cvnd/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/cvn:Z cvn/node "./cvnd testnet init-files --v 4 -o /cvn --keyring-backend=test --starting-ip-address 192.167.10.2"; fi
	docker-compose up -d

# Stop testnet
localnet-stop:
	docker-compose down

# Clean testnet
localnet-clean:
	docker-compose down
	sudo rm -rf build/*

 # Reset testnet
localnet-unsafe-reset:
	docker-compose down
ifeq ($(OS),Windows_NT)
	@docker run --rm -v $(CURDIR)\build\node0\cvnd:/cvn\Z cvn/node "./cvnd tendermint unsafe-reset-all --home=/cvn"
	@docker run --rm -v $(CURDIR)\build\node1\cvnd:/cvn\Z cvn/node "./cvnd tendermint unsafe-reset-all --home=/cvn"
	@docker run --rm -v $(CURDIR)\build\node2\cvnd:/cvn\Z cvn/node "./cvnd tendermint unsafe-reset-all --home=/cvn"
	@docker run --rm -v $(CURDIR)\build\node3\cvnd:/cvn\Z cvn/node "./cvnd tendermint unsafe-reset-all --home=/cvn"
else
	@docker run --rm -v $(CURDIR)/build/node0/cvnd:/cvn:Z cvn/node "./cvnd tendermint unsafe-reset-all --home=/cvn"
	@docker run --rm -v $(CURDIR)/build/node1/cvnd:/cvn:Z cvn/node "./cvnd tendermint unsafe-reset-all --home=/cvn"
	@docker run --rm -v $(CURDIR)/build/node2/cvnd:/cvn:Z cvn/node "./cvnd tendermint unsafe-reset-all --home=/cvn"
	@docker run --rm -v $(CURDIR)/build/node3/cvnd:/cvn:Z cvn/node "./cvnd tendermint unsafe-reset-all --home=/cvn"
endif

# Clean testnet
localnet-show-logstream:
	docker-compose logs --tail=1000 -f

.PHONY: localnet-build localnet-start localnet-stop localnet-clean localnet-unsafe-reset localnet-show-logstream
