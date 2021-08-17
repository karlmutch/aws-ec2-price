GOROOT := $(shell go env GOROOT)
GOPATH := $(shell go env GOPATH)

all: build docker-build

build:
	@go build -o $(GOPATH)/bin/aws-ec2-price ./cmd

docker-build:
	@docker build -t aws-ec2-price:latest -f dockerfiles/Dockerfile .
