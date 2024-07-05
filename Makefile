PROJECT=reddio

default: build

build:
	make -C juno rustdeps
	go build -v -o $(PROJECT) ./cmd/node/main.go ./cmd/node/testrequest.go

reset:
	@rm -r yu

clean:
	rm -f $(PROJECT)