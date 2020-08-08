VERSION=1.0

# helper to run go from project root
GO := go

# generate harbor cient files from swagger config
define install_provider
    mkdir -p ~/.terraform.d/plugins/linux_amd64/
    cp ./dist/terraform-provider-harbor_linux_amd64/terraform-provider-harbor_v* ~/.terraform.d/plugins/linux_amd64/
    chmod +x ~/.terraform.d/plugins/linux_amd64/terraform-provider-harbor_v*
endef

TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
MAKEFLAGS += --silent

default: build

generate:
	scripts/build-00-generate-client.sh

compile:
	# scripts/build-10-compile.sh
	goreleaser build --rm-dist --snapshot

install:
	$(call install_provider)

build: generate test compile

fmt:
	echo "==> Formatting files with fmt..."
	gofmt -w -s $(GOFMT_FILES)

test: goLint scriptsLint vet
	go test $(TEST)



fmtcheck:
	scripts/build-03-go-gofmtcheck.sh

vet:
	echo "==> Checking code with vet..."
	go vet ./...

goLint:
	scripts/build-03-go-gofmtcheck.sh
	scripts/build-04-go-errorchecks.sh
	scripts/build-05-go-golint.sh

gosec:
	echo "==> Checking code with gosec..."
	# TODO Remove unused files from generated sources !!!
	gosec -exclude-dir=gen/harborctl/client/scanners ./...

scriptsLint:
	echo "==> Checking scripts with shellcheck..."
	shellcheck scripts/*.sh

e2e_prepare:
	scripts/tst-00-prepare-kind.sh


e2e_prepare_harbor_v1:
	scripts/tst-01-prepare-harbor.sh "172-17-0-1.sslip.io" "1.3.2"

e2e_prepare_harbor_v2:
	scripts/tst-01-prepare-harbor.sh "172-17-0-1.sslip.io" "1.4.0"

e2e_prepare_harbor_v2_1:
	scripts/tst-01-prepare-harbor.sh "172-17-0-1.sslip.io" "1.4.1"

e2e_clean_cluster:
	kind delete cluster || true

e2e_clean_harbor:
	helm delete tf-harbor-test -n harbor
	sleep 10

e2e_test_v2:
	scripts/tst-15-execute-go-acc.sh "/api/v2.0"

e2e_test_v1:
	scripts/tst-15-execute-go-acc.sh "/api"

e2e_test_classic:
	scripts/tst-15-execute-classic-acc.sh "/api/v2.0"

e2e_full_run: e2e_clean_cluster e2e_prepare e2e_prepare_harbor_v2 e2e_test_v2 e2e_clean_harbor e2e_prepare_harbor_v1 e2e_test_v1 e2e_clean_cluster
# e2e_prepare e2e_prepare_harbor_v1 e2e_test e2e_cleanup
spellingCheck:
	mdspell '**/*.md' '!**/node_modules/**/*.md'

.PHONY: default install
