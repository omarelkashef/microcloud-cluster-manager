GOMIN=1.23.4
CGO_ENABLED?=0 # create statically linked binary
GOOS?=linux
GO_BIN?=app # name of the output application binary
GO?=go # name of the go binary
GOFLAGS?=-ldflags=-w -ldflags=-s -a -buildvcs # remove debug info, strip symbol table, force packages rebuild, includes version control metadata in binary
GO_UI_FOLDER?=internal/app/management-api/api/v1/static

# export all variables defined as environment variables
.EXPORT_ALL_VARIABLES:

.PHONY: default
default: all

# ==============================================================================
# Static code linting utility targets.

.PHONY: check
check: check-static

.PHONY: check-static
check-static:
ifeq ($(shell command -v golangci-lint 2> /dev/null),)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin
endif
	golangci-lint run --timeout 5m

# ==============================================================================
# Go module utility targets.

.PHONY: update-gomod
update-gomod:
	go get -t -v -d -u ./...
	go mod tidy -go=$(GOMIN)

.PHONY: tidy-gomod
tidy-gomod:
	go mod tidy -go=$(GOMIN)

# ====================================================================
# Local dev cluster utility targets. (k8s, kustomize, kind, skaffold)

KIND_CLUSTER := dev-cluster

.PHONY: start-cluster
start-cluster:
	@if ! kind get clusters | grep -q "$(KIND_CLUSTER)"; then \
		echo "Cluster '$(KIND_CLUSTER)' does not exist. Creating..."; \
		kind create cluster \
			--image kindest/node:v1.31.0 \
			--name $(KIND_CLUSTER) \
			--config deployment/k8s/kind/kind-config.yaml; \
		kubectl config set-context --current --namespace=default; \
	else \
		echo "Cluster '$(KIND_CLUSTER)' already exists."; \
	fi

.PHONY: delete-cluster
delete-cluster:
	kind delete cluster --name $(KIND_CLUSTER)

.PHONY: dev-k8s-configs
dev-k8s-configs:
	kubectl kustomize deployment/k8s/dev

.PHONY: dev-k8s-deploy
dev-k8s-deploy:
	skaffold dev --no-prune=false -p docker

.PHONY: debug-k8s-deploy
debug-k8s-deploy:
	skaffold dev --no-prune=false -p debug

.PHONY: rock-k8s-deploy
rock-k8s-deploy:
	skaffold dev --no-prune=false --cache-artifacts=false -p rock

# unfortunately necessary as skaffold does not automatically remove images after removing k8s cluster objects
.PHONY: clean-dev
clean-dev:
	docker container prune -f
	docker images -f "dangling=true" -q | xargs -r docker rmi
	docker images --filter=reference='lxd-cluster-manager:*' -q | xargs -I {} docker rmi {} -f

.PHONY: dev
dev: start-cluster dev-k8s-deploy

.PHONY: debug
debug: start-cluster debug-k8s-deploy

.PHONY: dev-rock
dev-rock: start-cluster rock-k8s-deploy

.PHONY: nuke
nuke: clean-dev delete-cluster

# ====================================================================
# UI utilities
.PHONY: ui
ui: 
	cd ui && dotrun

# ====================================================================
# dev database utilities

.PHONY: migrate-db
migrate-db:
	go run cmd/admin/main.go

# ====================================================================
# test utilities

# to ensure that all pods are ready before running tests, we check the liveliness of the pods
# rollout restart seems to break k8s portforwarding, here we make a request to the server to ensure it is up as well as reset the portforwarding
.PHONY: switch-test-mode
switch-test-mode:
	kubectl patch configmap config --patch '{"data":{"TEST_MODE":"$(IS_ON)"}}'
	kubectl rollout restart deployment/management-api-depl
	kubectl rollout status deployment/management-api-depl --timeout=300s
	@{ curl --insecure https://localhost:9000 > /dev/null 2>&1 || true; } 2>/dev/null

# Need to set TEST_MODE to true in the management-api deployment so we can by pass oidc authentication
.PHONY: test-e2e
test-e2e: 
	$(MAKE) switch-test-mode IS_ON=true
	go test -count=1 -v ./test/e2e
	$(MAKE) switch-test-mode IS_ON=false

.PHONY: test-ui-e2e
test-ui-e2e:
	cd ui && npx playwright test

# ====================================================================
# production build utilities for rockcraft

.PHONY: build-ui
build-ui:
	cd ui && yarn install --frozen=lockfile
	rm -rf ui/build
	cd ui && yarn build

.PHONY: copy-ui
copy-ui:
	rm -rf $(GO_UI_FOLDER)
	mkdir -p $(GO_UI_FOLDER)
	cp -r ui/build/ui $(GO_UI_FOLDER)

# create a binary "app" located in project root
.PHONY: build
build: build-ui copy-ui
	$(GO) build -C cmd -o $(GO_BIN) ./