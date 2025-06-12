PROJECT=reddio

default: build

build:
	go build -v -o ./$(PROJECT) ./cmd/node/main.go ./cmd/node/testrequest.go

## for local dev

build_transfer_test_no_race:
	go build -v -o ./bin/transfer_test ./test/cmd/transfer/main.go

build_uniswap_test_no_race:
	go build -v -o ./bin/uniswap_test ./test/cmd/uniswap/main.go

build_uniswap_benchmark_test:
	go build -v -o ./bin/uniswap_benchmark_test ./test/cmd/uniswap_benchmark/main.go

build_transfer_erc20_test_no_race:
	go build -v -o ./bin/transfer_erc20_test ./test/cmd/erc20/main.go

## for ci

build_transfer_test_race:
	go build -race -v -o ./bin/transfer_test ./test/cmd/transfer/main.go

build_uniswap_test_race:
	go build -race -v -o ./bin/uniswap_test ./test/cmd/uniswap/main.go

build_transfer_erc20_test_race:
	go build -race -v -o ./bin/transfer_erc20_test ./test/cmd/erc20/main.go

ci_parallel_transfer_test: reset
	./bin/transfer_test --parallel=true

ci_serial_transfer_test: reset
	./bin/transfer_test --parallel=false

ci_parallel_uniswap_test: reset
	./bin/uniswap_test --parallel=true

ci_serial_uniswap_test: reset
	./bin/uniswap_test --parallel=false

ci_parallel_transfer_erc20_test: reset
	./bin/transfer_erc20_test --parallel=true

ci_serial_transfer_erc20_test: reset
	./bin/transfer_erc20_test --parallel=FALSE

ci_transfer_client_test: reset
	./test/cmd/transfer/client_test.sh

## for local benchmark

build_benchmark_test:
	go build -v -o ./bin/benchmark_test ./test/cmd/benchmark/main.go

build_staging_benchmark_test:
	go build -v -o ./bin/staging_benchmark_test ./test/cmd/staging_benchmark/main.go

reset:
	@if [ -d "yu" ]; then \
		echo "Deleting 'yu' directory..."; \
		rm -rf yu; \
	fi
	@if [ -d "reddio_db" ]; then \
		echo "Deleting 'reddio_db' directory..."; \
		rm -rf reddio_db; \
	fi

clean:
	rm -f $(PROJECT)

clean_tests:
	rm uniswap_test benchmark_test

clean_test_data:
	@if [ -d "test/tmp" ]; then \
		echo "Deleting 'test/tmp' directory..."; \
		rm -rf test/tmp; \
	fi

check-mod-tidy:
	@go mod tidy
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "Changes detected after running go mod tidy. Please run 'go mod tidy' locally and commit the changes."; \
		git status; \
		exit 1; \
	else \
		echo "No changes detected after running go mod tidy."; \
	fi


# Define the paths to the ABI and BIN files
ERC20T_ABI = test/contracts/ERC20T.abi
ERC20T_BIN = test/contracts/ERC20T.bin
TOKEN_ABI = test/contracts/Token.abi
TOKEN_BIN = test/contracts/Token.bin
WETH9_ABI = test/contracts/WETH9.abi
WETH9_BIN = test/contracts/WETH9.bin
UNISWAPV2FACTORY_ABI = test/contracts/UniswapV2Factory.abi
UNISWAPV2FACTORY_BIN = test/contracts/UniswapV2Factory.bin
UNISWAPV2ROUTER01_ABI = test/contracts/UniswapV2Router01.abi
UNISWAPV2ROUTER01_BIN = test/contracts/UniswapV2Router01.bin
CHILDBRIDGECOREFACET_ABI = bridge/contract/ChildBridgeCoreFacet.abi
PARENTBRIDGECOREFACET_ABI = bridge/contract/ParentBridgeCoreFacet.abi
UPWARDMESSAGEDISPATCHERFACET_ABI = bridge/contract/UpwardMessageDispatcherFacet.abi
DOWNWARDMESSAGEDISPATCHERFACET_ABI = bridge/contract/DownwardMessageDispatcherFacet.abi
ERC20TOKEN_ABI = bridge/test/bindings/ERC20Token.abi
ERC721TOKEN_ABI = bridge/test/bindings/ERC721Token.abi
ERC1155TOKEN_ABI = bridge/test/bindings/ERC1155Token.abi
PARENTTOKENMESSAGETRANSMITTERFACET_ABI = bridge/test/bindings/ParentTokenMessageTransmitterFacet.abi
CHILDTOKENMESSAGETRANSMITTERFACET_ABI = bridge/test/bindings/ChildTokenMessageTransmitterFacet.abi
CHILDBRIDGECOREFACET_TEST_ABI = bridge/test/bindings/ChildBridgeCoreFacet.abi

# Define the output paths for the generated Go files
ERC20T_GO = test/contracts/ERC20T.go
TOKEN_GO = test/contracts/Token.go
WETH9_GO = test/contracts/WETH9.go
UNISWAPV2FACTORY_GO = test/contracts/UniswapV2Factory.go
UNISWAPV2ROUTER01_GO = test/contracts/UniswapV2Router01.go
CHILDBRIDGECOREFACET_GO = bridge/contract/ChildBridgeCoreFacet.go
PARENTBRIDGECOREFACET_GO = bridge/contract/ParentBridgeCoreFacet.go
UPWARDMESSAGEDISPATCHERFACET_GO = bridge/contract/UpwardMessageDispatcherFacet.go
DOWNWARDMESSAGEDISPATCHERFACET_GO = bridge/contract/DownwardMessageDispatcherFacet.go
ERC20TOKEN_GO = bridge/test/bindings/ERC20Token.go
ERC721TOKEN_GO = bridge/test/bindings/ERC721Token.go
ERC1155TOKEN_GO = bridge/test/bindings/ERC1155Token.go
PARENTTOKENMESSAGETRANSMITTERFACET_GO = bridge/test/bindings/ParentTokenMessageTransmitterFacet.go
CHILDTOKENMESSAGETRANSMITTERFACET_GO = bridge/test/bindings/ChildTokenMessageTransmitterFacet.go
CHILDBRIDGECOREFACET_TEST_GO = bridge/test/bindings/ChildBridgeCoreFacet.go


# Define the package name
PKG = contracts

# Define the abigen command
ABIGEN = abigen

# Rule to generate all Go bindings
generate_bindings:
	$(ABIGEN) --abi $(ERC20T_ABI) --bin $(ERC20T_BIN) --pkg $(PKG) --type ERC20T --out $(ERC20T_GO)
	$(ABIGEN) --abi $(TOKEN_ABI) --bin $(TOKEN_BIN) --pkg $(PKG) --type Token --out $(TOKEN_GO)
	$(ABIGEN) --abi $(WETH9_ABI) --bin $(WETH9_BIN) --pkg $(PKG) --type WETH9 --out $(WETH9_GO)
	$(ABIGEN) --abi $(UNISWAPV2FACTORY_ABI) --bin $(UNISWAPV2FACTORY_BIN) --pkg $(PKG) --type UniswapV2Factory --out $(UNISWAPV2FACTORY_GO)
	$(ABIGEN) --abi $(UNISWAPV2ROUTER01_ABI) --bin $(UNISWAPV2ROUTER01_BIN) --pkg $(PKG) --type UniswapV2Router01 --out $(UNISWAPV2ROUTER01_GO)
	# Clean up generated files
clean_bindings:
	rm -f $(ERC20T_GO) $(TOKEN_GO) $(WETH9_GO) $(UNISWAPV2FACTORY_GO) $(UNISWAPV2ROUTER01_GO)

BRIDGE_PKG = contract

generate_bridge_bindings:
	$(ABIGEN) --abi $(CHILDBRIDGECOREFACET_ABI) --pkg $(BRIDGE_PKG) --type ChildBridgeCoreFacet --out $(CHILDBRIDGECOREFACET_GO)
	$(ABIGEN) --abi $(PARENTBRIDGECOREFACET_ABI) --pkg $(BRIDGE_PKG) --type ParentBridgeCoreFacet --out $(PARENTBRIDGECOREFACET_GO)
	$(ABIGEN) --abi $(UPWARDMESSAGEDISPATCHERFACET_ABI) --pkg $(BRIDGE_PKG) --type UpwardMessageDispatcherFacet --out $(UPWARDMESSAGEDISPATCHERFACET_GO)
	$(ABIGEN) --abi $(DOWNWARDMESSAGEDISPATCHERFACET_ABI) --pkg $(BRIDGE_PKG) --type DownwardMessageDispatcherFacet --out $(DOWNWARDMESSAGEDISPATCHERFACET_GO)


TEST_PKG = bindings
generate_intergration_test_bindings:
	$(ABIGEN) --abi $(ERC20TOKEN_ABI) --pkg $(TEST_PKG) --type ERC20Token --out $(ERC20TOKEN_GO)
	$(ABIGEN) --abi $(ERC721TOKEN_ABI) --pkg $(TEST_PKG) --type ERC721Token --out $(ERC721TOKEN_GO)
	$(ABIGEN) --abi $(ERC1155TOKEN_ABI) --pkg $(TEST_PKG) --type ERC1155Token --out $(ERC1155TOKEN_GO)
	$(ABIGEN) --abi $(PARENTTOKENMESSAGETRANSMITTERFACET_ABI) --pkg $(TEST_PKG) --type ParentTokenMessageTransmitterFacet --out $(PARENTTOKENMESSAGETRANSMITTERFACET_GO)
	$(ABIGEN) --abi $(CHILDTOKENMESSAGETRANSMITTERFACET_ABI) --pkg $(TEST_PKG) --type ChildTokenMessageTransmitterFacet --out $(CHILDTOKENMESSAGETRANSMITTERFACET_GO)
	$(ABIGEN) --abi $(CHILDBRIDGECOREFACET_TEST_ABI) --pkg $(TEST_PKG) --type ChildBridgeCoreFacet --out $(CHILDBRIDGECOREFACET_TEST_GO)