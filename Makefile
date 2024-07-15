PROJECT=reddio

default: build

build:
	go build -v -o $(PROJECT) ./cmd/node/main.go ./cmd/node/testrequest.go

reset:
	@rm -r yu reddio_db

transfer_test:
	go build -v -o transfer_test ./test/trasnfer/main.go

clean:
	rm -f $(PROJECT)