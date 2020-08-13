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
	cd ./tools && go run mage.go -v build:generateHarborGoClient

compile:
	# scripts/build-10-compile.sh
	cd ./tools && go run mage.go -v build:goRelease

install:
	cd ./tools && go run mage.go -v build:TerraformInstallProvider

build: generate test compile

fmt:
	echo "==> Formatting files with fmt..."
	cd ./tools && go run mage.go -v build:fmt

test: goLint scriptsLint vet
	go test $(TEST)

fmtcheck:
	cd ./tools && go run mage.go -v build:fmt

vet:
	echo "==> Checking code with vet..."
	go vet ./...

goLint:
	fmt
	scripts/build-04-go-errorchecks.sh
	cd ./tools && go run mage.go -v lint || true

gosec:
	echo "==> Checking code with gosec..."
	# TODO Remove unused files from generated sources !!!
	gosec -exclude-dir=gen/harborctl/client/scanners ./...

scriptsLint:
	echo "==> Checking scripts with shellcheck..."
	shellcheck scripts/*.sh

e2e_prepare:
	cd ./tools && go run mage.go -v kind:recreate

e2e_prepare_harbor_v1:
	cd ./tools && HARBOR_HELM_CHART_VERSION=1.3.2 go run mage.go -v testArtefacts:deploy

e2e_prepare_harbor_v2:
	cd ./tools && HARBOR_HELM_CHART_VERSION=1.4.0 go run mage.go -v testArtefacts:deploy

e2e_prepare_harbor_v2_1:
	cd ./tools && HARBOR_HELM_CHART_VERSION=1.4.1 go run mage.go -v testArtefacts:deploy

e2e_clean_cluster:
	kind delete cluster || true

e2e_clean_harbor:
	cd ./tools && go run mage.go -v testArtefacts:delete
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
