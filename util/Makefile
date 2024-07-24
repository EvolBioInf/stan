nw = $(shell which noweb)

all : util

util: util.go
	go build util.go
util.go: util.org
	if [ "$(nw)" != "" ]; then\
		awk -f ../scripts/preTangle.awk util.org | bash ../scripts/org2nw | notangle -Rutil.go | gofmt > util.go;\
	fi
.PHONY: test
test: util_test.go util.go
	go test -v
util_test.go:
	if [ "$(nw)" != "" ]; then\
		awk -f ../scripts/preTangle.awk util.org | bash ../scripts/org2nw | notangle -Rutil_test.go | gofmt > util_test.go;\
	fi
tangle: util_test.go util.go

clean:
	rm -f *.go
