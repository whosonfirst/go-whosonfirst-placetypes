prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/whosonfirst/go-whosonfirst-placetypes; then rm -rf src/github.com/whosonfirst/go-whosonfirst-placetypes; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-placetypes
	cp *.go src/github.com/whosonfirst/go-whosonfirst-placetypes/

deps:   self

fmt:
	go fmt *.go

