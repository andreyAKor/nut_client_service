GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

build:
	@# for current arch. system
	@go build -ldflags="-s -w" -o '$(GOBIN)/nut_client_service' ./cmd/nut_client_service/main.go || exit
	@# for MIPS arch. system on Onion Omega2/Omega2+
	@GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -ldflags="-s -w" -o '$(GOBIN)/mips/nut_client_service' ./cmd/nut_client_service/main.go || exit

run:
	@go build -o '$(GOBIN)/nut_client_service' ./cmd/nut_client_service/main.go
	@'$(GOBIN)/nut_client_service' --config='$(GOBASE)/configs/nut_client_service.yml'

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
