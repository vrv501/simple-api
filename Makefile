#!/usr/bin/make

export GOPRIVATE =  # comma seperated values without quotes & spaces
GO_TEST_CMD = CGO_ENABLED=1 go test -race
GO_BUILD_CMD = CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

all: lint test build
.PHONY: all

build: generate tidy
	$(GO_BUILD_CMD) -o bin/app -a -v cmd/server/main.go
.PHONY: build

lint: tidy .golangci.yml
	@go vet ./...
	@golangci-lint --version
	golangci-lint run --config .golangci.yml ./...
.PHONY: lint

tidy:
	@go mod tidy
	@go mod download
.PHONY: tidy

degenerate:
	@find . -type f -name '*_mock.go' -delete
.PHONY: degenerate

generate: degenerate
	@go generate -x ./...
.PHONY: generate

test: generate tidy
	$(GO_TEST_CMD) -v ./...
.PHONY: test

coverage: generate tidy
	@mkdir -p coverage
	$(GO_TEST_CMD) -covermode=atomic -coverpkg=./... -coverprofile=coverage/coverage.out -v ./...
.PHONY: coverage

clean:
	@rm -rf coverage bin vendor
.PHONY: clean

run:
	@go run ./cmd/server/main.go
.PHONY: run

clean-build: degenerate clean build
.PHONY: clean-build
