CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/whosonfirst/go-storagemaster; then rm -rf src/github.com/whosonfirst/go-storagemaster; fi
	mkdir -p src/github.com/whosonfirst/go-storagemaster/provider
	cp storagemaster.go src/github.com/whosonfirst/go-storagemaster/
	cp provider/*.go src/github.com/whosonfirst/go-storagemaster/provider/
	cp -r vendor/src/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/aws/aws-sdk-go"

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor/src; then rm -rf vendor/src; fi
	cp -r src vendor/src
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go
	go fmt provider/*.go
	go fmt *.go

bin: 	rmdeps self
	@GOPATH=$(GOPATH) go build -o bin/storagemaster cmd/storagemaster.go
