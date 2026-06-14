# DCC OPA Bridge Makefile

GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=opa-dcc-extension

all: test build

build:
	@echo "Building OPA DCC Extension..."
	cd src && $(GOBUILD) -v -o ../bin/$(BINARY_NAME) .

test:
	@echo "Running Go Unit Tests..."
	$(GOTEST) -v ./src/...

test-integration:
	@echo "Running Python Logic Verification..."
	python3 tests/verify_opa.py

build-opa-extension:
	@echo "Building OPA with custom DCC built-in..."
	# Requires OPA build tool installed
	opa build -o dcc-policy.tar.gz ./src/

clean:
	rm -f bin/$(BINARY_NAME)
	rm -f dcc-policy.tar.gz

.PHONY: all build test test-integration build-opa-extension clean
