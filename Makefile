.PHONY: test
test:
	go test -v -failfast ./...

bench:
	go test -bench . -benchmem
