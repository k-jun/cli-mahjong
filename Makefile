.PHONY: test
test:
	go test -v -failfast ./...

test_integration:
	go test -v -failfast -tags=integration ./...

bench:
	go test -bench . -benchmem

debug:
	go test -v -tags=integration -run TestPlay1 ./...

coverage:
	go test -coverprofile=cover.out  ./...
	go tool cover -html=cover.out -o cover.html
	open cover.html
