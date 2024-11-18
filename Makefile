PROJECT=reddio

default: build

build:
	go build -v -o $(PROJECT) ./cmd/node/main.go ./cmd/node/testrequest.go

## for local dev

build_transfer_test_no_race:
	go build -v -o transfer_test ./test/cmd/transfer/main.go

build_uniswap_test_no_race:
	go build -v -o uniswap_test ./test/cmd/uniswap/main.go

build_uniswap_benchmark_test:
	go build -v -o uniswap_benchmark_test ./test/cmd/uniswap_benchmark/main.go

## for ci

build_transfer_test_race:
	go build -race -v -o transfer_test ./test/cmd/transfer/main.go

build_uniswap_test_race:
	go build -race -v -o uniswap_test ./test/cmd/uniswap/main.go

ci_parallel_transfer_test: reset
	./transfer_test --parallel=true

ci_serial_transfer_test: reset
	./transfer_test --parallel=false

ci_parallel_uniswap_test: reset
	./uniswap_test --parallel=true

ci_serial_uniswap_test: reset
	./uniswap_test --parallel=false

build_state_root_test:
	go build -v -o state_root_test ./test/cmd/state_root/main.go

state_root_test_gen:
	./state_root_test --action=gen

state_root_test_assert:
	./state_root_test --action=assert

## for local benchmark

build_benchmark_test:
	go build -v -o benchmark_test ./test/cmd/benchmark/main.go

parallel_benchmark_test:
	./benchmark_test --parallel=true --maxBlock=50 --qps=1000 --embedded=false

serial_benchmark_test:
	./benchmark_test --parallel=false --maxBlock=50 --qps=1000 --embedded=false

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

# Define the output paths for the generated Go files
ERC20T_GO = test/contracts/ERC20T.go
TOKEN_GO = test/contracts/Token.go
WETH9_GO = test/contracts/WETH9.go
UNISWAPV2FACTORY_GO = test/contracts/UniswapV2Factory.go
UNISWAPV2ROUTER01_GO = test/contracts/UniswapV2Router01.go

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