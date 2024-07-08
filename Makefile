PROJECT=reddio

default: build

build:
	go build -v -o $(PROJECT) ./cmd/node/main.go ./cmd/node/testrequest.go

reset:
	@rm -r yu reddio_db

clean:
	rm -f $(PROJECT)