VERSION ?= dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w -X github.com/urugus/baseline-cli/internal/cli.version=$(VERSION) -X github.com/urugus/baseline-cli/internal/cli.commit=$(COMMIT) -X github.com/urugus/baseline-cli/internal/cli.date=$(DATE)

.PHONY: build test fmt install clean release

build:
	go build -ldflags "$(LDFLAGS)" -o baseline ./cmd/baseline

test:
	go test ./...

fmt:
	gofmt -w .

install:
	go install -ldflags "$(LDFLAGS)" ./cmd/baseline

clean:
	rm -rf dist baseline

release: test
	mkdir -p dist
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/baseline-darwin-arm64 ./cmd/baseline
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/baseline-darwin-amd64 ./cmd/baseline
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/baseline-linux-arm64 ./cmd/baseline
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/baseline-linux-amd64 ./cmd/baseline
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/baseline-windows-amd64.exe ./cmd/baseline
	cd dist && shasum -a 256 * > checksums.txt
