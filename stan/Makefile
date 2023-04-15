version = $(shell git describe)
date = $(shell git log | grep Date | head -n 1 | sed -r 's/Date: +[A-Z][a-z]+ ([A-Z][a-z]+) ([0-9]+) [^ ]+ ([0-9]+) .+/\2_\1_\3/')

exe = stan
tests = test1.sh test2.sh test3.sh test4.sh test5.sh test6.sh test7.sh test8.sh test9.sh
all : $(exe)

$(exe): $(exe).go ../util/util.go
	go build -ldflags "-X github.com/evolbioinf/$(exe)/util.version=$(version) -X github.com/evolbioinf/$(exe)/util.date=$(date)" $(exe).go
$(exe).go: $(exe).org
	bash ../scripts/org2nw $(exe).org | notangle -R$(exe).go | gofmt > $(exe).go
../util/util.go : ../util/util.org
	make -C ../util
test: $(exe)_test.go $(tests)
	@go test
$(exe)_test.go: $(exe).org
	bash ../scripts/org2nw $(exe).org | notangle -R$(exe)_test.go | gofmt > $(exe)_test.go
$(tests): $(exe).org
	for test in $(tests); do \
		bash ../scripts/org2nw $(exe).org | notangle -R$$test > $$test; \
	done
clean:
	rm -f *.go $(tests)

