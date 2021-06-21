BINDIR ?= $(CURDIR)/bin
TMPDIR ?= $(CURDIR)/tmp
ARCH   ?= amd64

help:  ## display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: help build image all clean dev

test: ## test board
	go test ./...

dev: ## live reload development
	gin --path . --appPort 8080 --all --immediate --bin tmp/board run

build: ## build board
	mkdir -p $(BINDIR)
	CGO_ENABLED=0 go build -o ./bin/board

verify: test build ## tests and builds board

image: ## build docker image
	docker build -t lesterpig/board:latest .

clean: ## clean up created files
	rm -rf \
		$(BINDIR) \
		$(TMPDIR)

all: clean test build image ## runs test, build and image

test-coverage: ## Generate test coverage report
	mkdir -p $(TMPDIR)
	go test ./... --coverprofile $(TMPDIR)/outfile
	go tool cover -html=$(TMPDIR)/outfile

lint: ## Generate static analysis report
	goreportcard-cli -v
	golint ./...
