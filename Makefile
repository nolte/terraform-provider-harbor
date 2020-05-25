VERSION=1.0

# helper to run go from project root
GO := go

# generate harbor cient files from swagger config
define install_provider
	tar -zxf bin/terraform-provider-harbor_*_linux_amd64.tar.gz -C  ~/.terraform.d/plugins/linux_amd64/
    chmod +x ~/.terraform.d/plugins/linux_amd64/terraform-provider-harbor
endef

TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
MAKEFLAGS += --silent

default: build

generate:
	scripts/build-00-generate-client.sh

compile:
	scripts/build-10-compile.sh

install:
	$(call install_provider)

build: generate test compile

fmt:
	gofmt -w -s $(GOFMT_FILES)

test: fmtcheck vet
	go test $(TEST)

testacc: fmtcheck vet
	TF_ACC=1 go test -timeout 20m $(TEST) -v $(TESTARGS)

fmtcheck:
	lineCount=$(shell gofmt -l -s $(GOFMT_FILES) | wc -l | tr -d ' ') && exit $$lineCount

vet:
	go vet ./...

lint:
	golangci-lint run

check-scripts:
	shellcheck scripts/*.sh
	shellcheck scripts/test/bats/build/*.bats

e2e_prepare:
	scripts/tst-00-prepare-kind.sh
	scripts/tst-01-prepare-harbor.sh

e2e_cleanup:
	kind delete cluster

.PHONY: default install
