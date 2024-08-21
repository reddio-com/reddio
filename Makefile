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

transfer_test: reset build_transfer_test_race
	./transfer_test

parallel_transfer_test: reset
	./transfer_test --parallel=true

serial_transfer_test: reset
	./transfer_test --parallel=false

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