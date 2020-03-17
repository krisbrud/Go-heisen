# TODO:
#	runsimulator
#	runelevator
#	build?

.PHONY: run
run:
	go run src/main/main.go

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	gofmt -w src/**/*.go
