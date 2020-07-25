#!/usr/bin/env bash
# Bash3 Boilerplate. Copyright (c) 2014, kvz.io

set -o errexit
set -o pipefail
set -o nounset
# set -o xtrace

# Set magic variables for current file & dir
__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__root="$(cd "$(dirname "${__dir}")" && pwd)"

VERSION="${1:-v0.1.0}"
projectBase=$__root
STARTDIR=$(pwd)
if [ -d "${projectBase}/bin" ]; then
    echo "Remove old generated binary folder"
    # shellcheck disable=SC2115
    rm -rf "${projectBase}/bin"
fi

mkdir -p "${projectBase}/bin"

package_name=terraform-provider-harbor
platforms=(
    "darwin/amd64"
    "linux/amd64"
    "windows/amd64"
)
for platform in "${platforms[@]}"; do
    echo "Building terraform-provider-harbor_${VERSION} for ${platform}"

    # shellcheck disable=SC2206
    platform_split=(${platform//\// })

    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$package_name
    if [ "$GOOS" = "windows" ]; then
        output_name+='.exe'
    fi
    env CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" go build -o "${projectBase}/bin/${output_name}"
    tar -czvf "bin/terraform-provider-harbor_${VERSION}_${GOOS}_${GOARCH}.tar.gz" -C "${projectBase}/bin" "${output_name}"
    rm bin/$output_name
done
cd "${projectBase}/bin"
sha256sum -b ./* >SHA256SUMS
cd "${STARTDIR}"
