# メタ情報
NAME := GraphQLChat
MAIN := server/server.go
##Gitに登録されたら以下を使用
##VERSION := $(shell git describe --tags --abbrev=0)
##REVISION := $(shell git rev-parse --short HEAD)
VERSION := v1.0.0
REVISION := 0000001
LDFLAGS :=	-X 'main.version=$(VERSION)' \
			-X 'main.revision=$(REVISION)'
GOMOD := go.mod
GOMODEXISTS := $(shell ls | grep $(GOMOD))
SRCS := $(shell find . -type f -name '*.go')

# env 
export GO111MODULE=on
export PORT=8080
export REDIS_URL=localhost:6379

# Setup
setup:
ifeq ($(GOMODEXISTS),$(GOMOD))
else
	go mod init
	go get -u github.com/99designs/gqlgen
	go get -u github.com/go-redis/redis
	go get golang.org/x/tools/cmd/goimports
	go get github.com/golang/lint/golint
	go get github.com/Songmu/make2help/cmd/make2help
endif

## Run tests
test: 
	@go test $(SRCS)

## Lint
lint: setup
	@for pkg in $$(go list ./...); do \
		golint --set_exit_status $$pkg || exit $$?; \
	done

## vet
vet: setup
	@go tool vet $(SRCS)

## Format source codes
fmt: setup
	@goimports -w $(SRCS)

## check go source (goimpots, vet, lint)
check:
	@make fmt
	@echo "Done! fmt"
	@make vet 
	@echo "Done! vet"
	@make lint 
	@echo "Done! lint"

## build binaries ex. make bin/$(NAME)
build: $(MAIN)
	go build -ldflags "$(LDFLAGS)" -o bin/$(NAME) $<

run: $(MAIN)
	go run $<

## gqlgen script
gqlgen: setup
	go run scripts/gqlgen.go ${ARG}

## Show help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: setup test lint vet fmt build run help
