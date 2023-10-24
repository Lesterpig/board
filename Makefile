cluster-name := tcs-board
kubeconfig := tmp/kubeconfig.yaml
image-name := tcs-board

.PHONY: cluster delete-cluster image load-image deploy
.PHONY: secure lint fmt test coverage build run tidy clean

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

cluster: k3d ## Create a k3d cluster
	KUBECONFIG="$(kubeconfig)" $(K3D) cluster create $(cluster-name) -p "8081:80@loadbalancer" --agents 2 2> /dev/null | true

delete-cluster: k3d ## Delete the k3d cluster
	KUBECONFIG="$(kubeconfig)" $(K3D) cluster delete $(cluster-name)

image: ## Build a docker image
	docker build -t $(image-name):latest .

load-image: image cluster ## Load the locally built image into k3d
	KUBECONFIG="$(kubeconfig)" $(K3D) image import $(image-name):latest --cluster $(cluster-name)

deploy: cluster ## Apply manifests in examples/ to Kubernetes
	kubectl apply -k examples/ --kubeconfig="$(kubeconfig)"

secure: gosec ## Run gosec
	$(GOSEC) -terse ./...

lint: golangci-lint ## Generate static analysis report
	GOLANGCI_LINT_CACHE="$(PWD)/tmp/.cache" $(GOLANGCILINT) run

fmt: ## Run go fmt
	go fmt ./...

test: ## Run Go unit-tests
	go test ./...

coverage: ## Generate test coverage report
	mkdir -p tmp/
	go test ./... -coverprofile=tmp/coverage.out
	go tool cover -func=tmp/coverage.out

build: ## Build the binary to build/main
	CGO_ENABLED=0 go build -o build/main ./main.go
	$(RICE) append --exec build/main

run: ## Run the program
	go run ./... -f examples/board.yaml

tidy: ## Run go mod tidy
	go mod tidy

clean: ## Clean up local repository
	rm -rf bin/
	rm -rf build/
	rm -rf tmp/

## Build dependencies

LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool binaries
K3D ?= $(LOCALBIN)/k3d
GOSEC ?= $(LOCALBIN)/gosec
GOLANGCILINT ?= $(LOCALBIN)/golangci-lint
RICE ?= $(LOCALBIN)/rice

.PHONY: k3d
k3d: $(K3D) ## Download k3d
$(K3D): $(LOCALBIN)
	test -s $(LOCALBIN)/k3d || GOBIN=$(LOCALBIN) go install github.com/k3d-io/k3d/v5@latest

.PHONY: gosec
gosec: $(GOSEC) ## Download gosec
$(GOSEC): $(LOCALBIN)
	test -s $(LOCALBIN)/gosec || GOBIN=$(LOCALBIN) go install github.com/securego/gosec/v2/cmd/gosec@latest

.PHONY: golangci-lint
golangci-lint: $(GOLANGCILINT) ## Download golangci-lint
$(GOLANGCILINT): $(LOCALBIN)
	test -s $(LOCALBIN)/golangci-lint || GOBIN=$(LOCALBIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: rice
rice: $(RICE) ## Download rice
$(RICE): $(LOCALBIN)
	test -s $(LOCALBIN)/rice || GOBIN=$(LOCALBIN) go install github.com/GeertJohan/go.rice/rice@latest

