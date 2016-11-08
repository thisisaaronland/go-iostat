CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src/github.com/thisisaaronland/go-iostat; then rm -rf src/github.com/thisisaaronland/go-iostat; fi
	mkdir -p src/github.com/thisisaaronland/go-iostat
	cp iostat.go src/github.com/thisisaaronland/go-iostat/
	# cp -r vendor/src/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:	rmdeps

vendor-deps: deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor/src; then rm -rf vendor/src; fi
	cp -r src vendor/src
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt *.go
	go fmt cmd/*.go

bin: 	self
	@GOPATH=$(GOPATH) go build -o bin/iostat cmd/iostat.go
