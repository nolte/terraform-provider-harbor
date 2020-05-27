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

test: goLint scriptsLint vet
	go test $(TEST)

testaac:
	scripts/tst-15-execute-go-acc.sh

fmtcheck:
	scripts/build-03-go-gofmtcheck.sh

vet:
	go vet ./...

goLint:
	scripts/build-03-go-gofmtcheck.sh
	scripts/build-04-go-errorchecks.sh
	scripts/build-05-go-golint.sh

scriptsLint:
	echo "==> Checking scripts with shellcheck..."
	shellcheck scripts/*.sh
	shellcheck scripts/test/bats/build/*.bats

e2e_prepare:
	scripts/tst-00-prepare-kind.sh
	scripts/tst-01-prepare-harbor.sh

e2e_cleanup:
	kind delete cluster

.PHONY: default install
