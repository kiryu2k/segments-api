.PHONY: build
build:
	@go build -o ./bin/segments ./cmd/segments/main.go

.PHONY: run
run: build
	@./bin/segments -config ./configs/config.dev.yaml

.PHONY: test
test:
	@go test -v ./...

.PHONY: .install-linter
.install-linter:
	@[ -f ./bin/golangci-lint ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.53.3

.PHONY: lint
lint: .install-linter
	@./bin/golangci-lint run ./...

.PHONY: lint-fast
lint-fast: .install-linter
	@./bin/golangci-lint run ./... --fast

.DEFAULT_GOAL := run