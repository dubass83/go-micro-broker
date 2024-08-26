.PHONY: *

docker_up:
	limactl start docker

docker_down: 
	limactl stop docker

test:
	go test -v -cover -count=1 -short ./...

server:
	go run main.go

build:
	docker build -t go-micro-broker -f Dockerfile .