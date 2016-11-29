CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	rmdeps deps fmt bin

self:   prep
	if test -d src/github.com/whosonfirst/go-whosonfirst-placetypes; then rm -rf src/github.com/whosonfirst/go-whosonfirst-placetypes; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-placetypes
	cp *.go src/github.com/whosonfirst/go-whosonfirst-placetypes/
	cp -r spec src/github.com/whosonfirst/go-whosonfirst-placetypes/spec/

deps:   self

fmt:
	go fmt *.go

spec:
	@GOPATH=$(GOPATH) go run cmd/mk-spec.go > placetypes/spec.go

# @GOPATH=$(shell pwd) go get -u ""
