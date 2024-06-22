#!/usr/bin/make

export GOPRIVATE =  # comma seperated values without quotes & spaces
GO_TEST_CMD = CGO_ENABLED=1 go test -race
GO_BUILD_CMD = CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

all: lint test build
.PHONY: all

vendor: go.mod go.sum
	@go mod download
	@go mod vendor
.PHONY: vendor

build: vendor
	$(GO_BUILD_CMD) -o bin/app -a -v cmd/main.go
.PHONY: build

lint: vendor .golangci.yml
	@go vet ./...
	@golangci-lint --version
	golangci-lint run --config .golangci.yml ./...
.PHONY: lint

tidy:
	@go mod tidy
.PHONY: tidy

degenerate:
	@find . -type f -name '*_mock.go' -delete
.PHONY: degenerate

generate: degenerate vendor
	@go generate ./...
.PHONY: generate

test: vendor generate
	$(GO_TEST_CMD) -v ./...
.PHONY: test

coverage: vendor generate
	@mkdir -p coverage
	$(GO_TEST_CMD) -covermode=atomic -coverpkg=./... -coverprofile=coverage/coverage.out -v ./...
.PHONY: coverage

clean:
	@rm -rf coverage bin vendor
.PHONY: clean

run:
	@go run ./cmd/main.go
.PHONY: run

clean-build: degenerate clean build
.PHONY: clean-build
