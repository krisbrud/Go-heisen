# TODO:
#	runsimulator
#	runelevator
#	build?

.PHONY: simulators
simulators:
	gnome-terminal -- SimElevatorServer  --port 14100
	sleep 0.5
	gnome-terminal -- go run src/main/main.go --port 14100

	gnome-terminal -- SimElevatorServer  --port 14101
	sleep 0.5
	gnome-terminal -- go run src/main/main.go --port 14101

	gnome-terminal -- SimElevatorServer  --port 14102
	sleep 0.5
	gnome-terminal -- go run src/main/main.go --port 14102


.PHONY: run
run:
	go run src/main/main.go

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	gofmt -w src/**/*.go
