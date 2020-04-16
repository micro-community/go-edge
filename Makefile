
NAME=x-edge
IMAGE_NAME=docker.pkg.github.com/micro-community/$(NAME)
GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_TAG=$(shell git describe --abbrev=0 --tags --always --match "v*")
GIT_IMPORT=github.com/micro-community/x-edge/app
CGO_ENABLED=0
BUILD_DATE=$(shell date +%s)
LDFLAGS=-X $(GIT_IMPORT).GitCommit=$(GIT_COMMIT) -X $(GIT_IMPORT).GitTag=$(GIT_TAG) -X $(GIT_IMPORT).BuildDate=$(BUILD_DATE)
IMAGE_TAG=$(GIT_TAG)-$(GIT_COMMIT)
GOPATH:=$(shell go env GOPATH)

all: build

vendor:
	go mod vendor

proto:
	protoc --proto_path=${GOPATH}/src:. --micro_out=. --go_out=. proto/protocol/proto_contract.proto

build:
	go build -a -installsuffix cgo -ldflags "-w ${LDFLAGS}" -o $(NAME) ./example/*.go

buildw:
	go build -a -installsuffix cgo -ldflags "-w ${LDFLAGS}" -o $(NAME).exe ./example/*.go

.PHONY: test
test:
	go test -v ./... -cover

docker:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(IMAGE_NAME):latest
	docker push $(IMAGE_NAME):$(IMAGE_TAG)
	docker push $(IMAGE_NAME):latest

.PHONY: buildw build clean vet test docker proto

