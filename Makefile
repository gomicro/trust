SHELL = bash

APP := $(shell basename $(PWD) | tr '[:upper:]' '[:lower:]')

GIT_COMMIT_HASH ?= $(shell git rev-parse HEAD)
GIT_SHORT_COMMIT_HASH := $(shell git rev-parse --short HEAD)

.PHONY: all
all: test

.PHONY: clean
clean: ## Cleans out all generated items
	-@rm -rf coverage
	-@rm -f coverage.txt

.PHONY: coverage
coverage: ## Generates the code coverage from all the tests
	docker run -v $$PWD:/go$${PWD/$$GOPATH} --workdir /go$${PWD/$$GOPATH} gomicro/gocover

.PHONY: generate_key
generate_key: ## Generate a new key
	openssl genrsa -des3 -out testCA.key 2048

.PHONY: generate_pem
generate_pem: generate_key ## Generate the pem from a key
	openssl req -x509 -new -nodes -key testCA.key -sha256 -days 3650 -out testCA.pem

.PHONY: help
help:  ## Show This Help
	@for line in $$(cat Makefile | grep "##" | grep -v "grep" | sed  "s/:.*##/:/g" | sed "s/\ /!/g"); do verb=$$(echo $$line | cut -d ":" -f 1); desc=$$(echo $$line | cut -d ":" -f 2 | sed "s/!/\ /g"); printf "%-30s--%s\n" "$$verb" "$$desc"; done

.PHONY: linters
linters: ## Run all the linters
	docker run -v $$PWD:/go$${PWD/$$GOPATH} --workdir /go$${PWD/$$GOPATH} gomicro/golinters

.PHONY: test
test: unit_test ## Run all available tests

.PHONY: update_rootca
update_rootca: ## Update the root CA from the latest copy of centos
	@echo -en "package trust\n\n" > globalchain.go
	@echo -en "const globalPemCerts string = \`\n\n" >> globalchain.go
	@docker run centos /bin/bash -c 'cat /etc/ssl/certs/ca-bundle.crt' >> globalchain.go
	@echo "\`" >> globalchain.go

.PHONY: unit_test
unit_test: ## Run unit tests
	go test ./...
