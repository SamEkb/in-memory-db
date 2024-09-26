LOCAL_BIN:=$(CURDIR)/bin

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3

lint:
	PATH=$(LOCAL_BIN):$(PATH) golangci-lint run ./... --config .golangci.pipeline.yaml

test_coverage:
	go test ./... -coverprofile=coverage.out