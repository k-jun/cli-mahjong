.PHONY: test
test:
	go test -v -failfast ./...

test_integration:
	go test -v -failfast -tags=integration ./...

bench:
	go test -bench . -benchmem

debug:
	go test -v -tags=integration -run TestPlay1 ./...
