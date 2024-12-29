.PHONY: *

docker_up:
	limactl start docker

docker_down: 
	limactl stop docker

docker_build:
	docker build -t go-micro-broker -f Dockerfile .

docker_build_simple: build
	docker build -t go-micro-broker -f Dockerfile.simple .

test:
	go test -v -cover -count=1 -short ./...

server:
	go run main.go

build:
	go build -o main main.go

proto:
	rm -f pb/*
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
      --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
      proto/*.proto