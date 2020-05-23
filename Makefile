VERSION=1.0

# helper to run go from project root
GO := go

# generate harbor cient files from swagger config
define install_provider
    mv terraform-provider-harbor_v1.0_linux_amd64 ~/.terraform.d/plugins/linux_amd64/terraform-provider-harbor
    chmod +x ~/.terraform.d/plugins/linux_amd64/terraform-provider-harbor
endef
# run build command
define building_provider
	echo "Building terraform-provider-harbor_${VERSION}_linux_amd64..."
	env GOOS=linux GOARCH=amd64 $(GO) build -o terraform-provider-harbor_v${VERSION}_linux_amd64 .
endef

default: build

generate:
	scripts/generate-client.sh

install:
	$(call install_provider)

build: generate build_cleanup
	$(call building_provider,build)

compile:
	$(call building_provider,install)


build_cleanup:
	rm -f ./terraform-provider-harbor_*

.PHONY: default generate build build_cleanup install