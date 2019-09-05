GO    := go
pkgs   = $(shell $(GO) list ./... | grep -v /vendor/)

all: test test-report

test:
	@echo ">> running tests"
	@$(GO) test $(pkgs) --cover

test-report:
	@echo ">> running tests"
	@$(GO) test $(pkgs) -coverprofile=coverage.txt && go tool cover -html=coverage.txt
