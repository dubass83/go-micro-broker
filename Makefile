.PHONY: *

docker_up:
	limactl start docker

docker_down: 
	limactl stop docker

docker_build:
	docker build -t go-micro-broker -f Dockerfile .

test:
	go test -v -cover -count=1 -short ./...

server:
	go run main.go
