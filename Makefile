PROJECT=reddio

default: build

build:
	go build -v -o $(PROJECT) ./cmd/node/main.go ./cmd/node/testrequest.go

build_transfer_test_race:
	go build -race -v -o transfer_test ./test/cmd/transfer/main.go

reset:
	@if [ -d "yu" ]; then \
		echo "Deleting 'yu' directory..."; \
		rm -rf yu; \
	fi
	@if [ -d "reddio_db" ]; then \
		echo "Deleting 'reddio_db' directory..."; \
		rm -rf reddio_db; \
	fi

benchmark_test: reset
	go run ./test/cmd/benchmark/main.go

parallel_benchmark_test: reset
	./benchmark_test --parallel=true --maxBlock=50 --qps=1000

serial_benchmark_test: reset
	./benchmark_test --parallel=false --maxBlock=50 --qps=1000

build_benchmark_test: reset
	go build -v -o benchmark_test ./test/cmd/benchmark/main.go

parallel_transfer_test: reset
	./transfer_test --parallel=true

serial_transfer_test: reset
	./transfer_test --parallel=false

build_uniswap_test: reset
	go build -race -v -o uniswap_test ./test/cmd/uniswap/main.go

parallel_uniswap_test: reset build_uniswap_test
	./uniswap_test --parallel=true

serial_uniswap_test: reset build_uniswap_test
	./uniswap_test --parallel=false


clean:
	rm -f $(PROJECT)

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