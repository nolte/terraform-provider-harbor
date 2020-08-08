#!/usr/bin/env bash

echo "==> Checking for unchecked errors..."

# shellcheck disable=SC2230
if ! which errcheck > /dev/null; then
    cd ./tools
    echo "==> Installing errcheck..."
    go get github.com/kisielk/errcheck
    cd ..
fi

# shellcheck disable=SC1126,SC2046
err_files=$(errcheck -ignoretests \
                     -ignore 'github.com/hashicorp/terraform/helper/schema:Set' \
                     -ignore 'bytes:.*' \
                     -ignore 'io:Close|Write' \
                     $(go list ./...| grep -v /vendor/))

if [[ -n "${err_files}" ]]; then
    echo 'Unchecked errors found in the following places:'
    echo "${err_files}"
    echo "Please handle returned errors. You can check directly with \`make errcheck\`"
    exit 1
fi

exit 0
