GOCMD=go
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
BINARY_NAME=ankiadder.bin
VERSION?=0.0.1
SERVICE_PORT?=3000

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

all: help

.PHONY: build
build: ## build app
	$(GOCMD) get ./...
	mkdir -p out/bin
	rm -rf out/bin/*
	$(GOCMD) build -o out/bin/$(BINARY_NAME) .


.PHONY: clean
clean: ## remove out folder and clean project
	rm -rf out
	$(GOCMD) clean

.PHONY: serve
serve: clean build ## clean and build project and finally run app
	./out/bin/$(BINARY_NAME)

.PHONY: run
run: ## run app
	$(GOCMD) run .

.PHONY: help
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)