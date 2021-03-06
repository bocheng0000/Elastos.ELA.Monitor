VERSION := $(shell git describe --abbrev=4 --dirty --always --tags)
BUILD = go build -ldflags "-X main.Version=$(VERSION) -X 'main.GoVersion=`go version`'" #-race

all:
	$(BUILD) -o monitor main.go

format:
	go fmt ./*

clean:
	rm -rf *.8 *.o *.out *.6