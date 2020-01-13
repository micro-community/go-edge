
GOPATH:=$(shell go env GOPATH)

proto:
	protoc --proto_path=${GOPATH}/src:. --micro_out=. --go_out=. proto/protocol/proto_contract.proto

build:
	go build -o x-edge example/main.go

.PHONY: test
test:
	go test -v ./... -cover

docker:
	docker build . -t x-edge:latest

.PHONY: docker build proto