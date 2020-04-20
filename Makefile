.PHONY: singlesim
singlesim:
	gnome-terminal -- SimElevatorServer  --port 14100
	sleep 0.5
	gnome-terminal -- go run src/main/main.go --port 14100

.PHONY: simulators
simulators:
	gnome-terminal --geometry=64x25+0+0 -- /bin/sh -c 'echo elev1; SimElevatorServer  --port=14101'
	sleep 0.5
	gnome-terminal --geometry=64x25+0+500 -- go run -race src/main/main.go --port=14101 --id=elev1 || read

	gnome-terminal --geometry=64x25+600+0 -- /bin/sh -c 'echo elev2; SimElevatorServer  --port=14102'
	sleep 0.5
	gnome-terminal --geometry=64x25+600+500 -- go run -race src/main/main.go --port=14102 --id=elev2 || read

	# Third elevator, currently ignored
	# gnome-terminal --geometry=64x25+1200+0 -- /bin/sh -c 'echo elev3; SimElevatorServer  --port=14103'
	# sleep 0.5
	# gnome-terminal --geometry=64x25+1200+500 -- go run src/main/main.go --port=14103 --id=elev3

.PHONY: packetloss
packetloss:
	sudo iptables -A INPUT -p udp --dport 44232 -m statistic --mode random --probability 0.2 -j DROP
	sudo iptables -A INPUT -p udp --dport 44233 -m statistic --mode random --probability 0.2 -j DROP
	sudo iptables -A INPUT -p udp --sport 44232 -m statistic --mode random --probability 0.2 -j DROP
	sudo iptables -A INPUT -p udp --sport 44233 -m statistic --mode random --probability 0.2 -j DROP	

.PHONY: disconnectnetwork
disconnectnetwork:
	sudo iptables -A INPUT -p udp --dport 44232 -m statistic --mode random --probability 1.0 -j DROP
	sudo iptables -A INPUT -p udp --dport 44233 -m statistic --mode random --probability 1.0 -j DROP
	sudo iptables -A INPUT -p udp --sport 44232 -m statistic --mode random --probability 1.0 -j DROP
	sudo iptables -A INPUT -p udp --sport 44233 -m statistic --mode random --probability 1.0 -j DROP	

.PHONY: normalnetwork
normalnetwork:
	sudo iptables -F

.PHONY: run
run:
	go run src/main/main.go

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	gofmt -w src/**/*.go
