.PHONY: vendor
vendor:
	go mod tidy && go mod vendor

# copy .env.example so we can adjust without affecting git file change
# actual server run does not required .env file to be exist.
.PHONY: env
env:
	cp .env.example .env

.PHONY: server/start
server/start:
	docker-compose -f docker-compose.yml up -d

.PHONY: server/restart
server/restart:
	docker-compose -f docker-compose.yml restart db-migration myapp

.PHONY: server/stop
server/stop:
	docker-compose -f docker-compose.yml down

################
# LINT AND STATIC ANALYSIS
################

GOLANGCI_VERSION=1.55.2
GOLANGCI_CHECK := $(shell golangci-lint version 2> /dev/null)

.PHONY: lint
lint:
# if golangci-lint failed on MacOS Ventura, try: brew install diffutils
ifndef GOLANGCI_CHECK
	# install by go get from source
	# go install github.com/golangci/golangci-lint/cmd/golangci-lint@v$(GOLANGCI_VERSION)
	# install by compiled binary per platform
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v$(GOLANGCI_VERSION)
endif
	golangci-lint run -c .golangci.yml ./...

################
# UNIT TEST
################

COVERAGE_OUTPUT=coverage.out
COVERAGE_OUTPUT_HTML=coverage.html
TEST_LIST := $(shell go list ./... | grep -v mock | grep -v /vendor/)
TEST_FLAGS=-failfast -v -race -trimpath -mod=readonly -coverprofile=$(COVERAGE_OUTPUT) -coverpkg=./...

.PHONY: test
test:
	go test $(TEST_LIST) $(TEST_FLAGS)
	go tool cover -html=$(COVERAGE_OUTPUT) -o $(COVERAGE_OUTPUT_HTML)
