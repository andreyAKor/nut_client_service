GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

build:
	@go build -ldflags="-s -w" -o '$(GOBIN)/nut_parser' ./cmd/nut_parser/main.go || exit
	@GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -ldflags="-s -w" -o '$(GOBIN)/nut_parser_mips' ./cmd/nut_parser/main.go || exit

run:
	@go build -o '$(GOBIN)/nut_parser' ./cmd/nut_parser/main.go
	@'$(GOBIN)/nut_parser' --config='$(GOBASE)/configs/nut_parser.yml'

up:
	@docker-compose up -d --build

down:
	@docker-compose down

test:
	@go test -v -count=1 -race -timeout=60s ./...

install-deps: deps
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint && go mod vendor && go mod verify

lint: install-deps
	@golangci-lint run ./...

deps:
	@go mod tidy && go mod vendor && go mod verify

install:
	@go mod download

generate:
	@go generate ./...

.PHONY: build
