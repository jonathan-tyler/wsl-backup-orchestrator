GO ?= go
DIST_DIR ?= dist
BINARY ?= backup
CMD_PATH ?= ./cmd/backup

.PHONY: all test test-unit test-integration build release clean

all: test build

test:
	$(GO) test ./...

test-unit:
	$(GO) test ./internal/...

test-integration:
	$(GO) test ./tests/integration/...

build:
	$(GO) build ./...

release: clean
	mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 $(GO) build -o $(DIST_DIR)/$(BINARY) $(CMD_PATH)

clean:
	rm -rf $(DIST_DIR)