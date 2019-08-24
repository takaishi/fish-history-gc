TEST ?= $(shell go list ./... | grep -v -e vendor -e keys -e tmp)

INFO_COLOR=\033[1;34m
RESET=\033[0m
BOLD=\033[1m

default: test lint

depsdev:
	go get -u golang.org/x/lint/golint

fmt:
	go fmt ./...

test: fmt
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Testing$(RESET)"
	go test -v $(TEST) -timeout=30s -parallel=4
	go test -race $(TEST)

lint: depsdev fmt
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Linting$(RESET)"
	golint -min_confidence 1.1 -set_exit_status $(TEST)
