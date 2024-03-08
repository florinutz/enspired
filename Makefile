VERSION ?= $(shell git describe --tags 2> /dev/null || echo v0)

GO = go
BINARY = room-parser

.PHONY: build clean test bench
.DEFAULT_GOAL := help

example: ## run the example
	${GO} run main.go rooms.txt

build: ## builds
	@${GO} build -ldflags "-X main.Version=${VERSION}" -o ${BINARY}
	@echo "${BINARY} built. Run it like this:\n\n\t./${BINARY} rooms.txt"

test: ## runs the unit tests
	${GO} test -v ./...

clean: ## go clean, then remove any previously built binary
	${GO} clean
	rm -f ${BINARY}

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

