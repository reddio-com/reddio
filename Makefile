PROJECT=reddio

default: build

build:
	go build -v -o $(PROJECT) ./cmd/node/main.go ./cmd/node/testrequest.go

reset:
	@if [ -d "yu" ]; then \
		echo "Deleting 'yu' directory..."; \
		rm -rf yu; \
	fi
	@if [ -d "reddio_db" ]; then \
		echo "Deleting 'reddio_db' directory..."; \
		rm -rf reddio_db; \
	fi

transfer_test: reset
	go run ./test/cmd/transfer/main.go

clean:
	rm -f $(PROJECT)