# TODO:
#	runsimulator
#	runelevator
#	build?

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	gofmt -w src/**/*.go
