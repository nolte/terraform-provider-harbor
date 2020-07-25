#!/usr/bin/env bash
# Bash3 Boilerplate. Copyright (c) 2014, kvz.io

set -o errexit
set -o pipefail
set -o nounset
# set -o xtrace

# Set magic variables for current file & dir
__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__root="$(cd "$(dirname "${__dir}")" && pwd)"

projectBase=$__root

GENERATED_SOURCES_TARGET="${projectBase}/gen/harborctl"
GENERATED_MERGED_SWAGGER="${projectBase}/gen/merged.json"

echo "==> Generate GoApi Client from Swagger Spec"

if [ -d "$GENERATED_SOURCES_TARGET" ]; then
    echo "Remove old generated sources"
    rm -rf "${GENERATED_SOURCES_TARGET}"
fi

if test -f "$GENERATED_MERGED_SWAGGER"; then
    echo "Remove old generated merged swagger conf"
    rm "${GENERATED_MERGED_SWAGGER}"
fi

mkdir -p "${projectBase}/gen"

# shellcheck disable=SC2002
cat "${projectBase}/scripts/swagger-specs/v2-swagger-original.json" | json-patch -p "${projectBase}/scripts/swagger-specs/patch.1.json" >"${GENERATED_MERGED_SWAGGER}"

mkdir -p "${GENERATED_SOURCES_TARGET}"

swagger generate client \
    -f "${GENERATED_MERGED_SWAGGER}" \
    --name=harbor \
    --target="${GENERATED_SOURCES_TARGET}" \
    --with-flatten=remove-unused
#--operation=PostProjects \
#--operation=GetProjects \
#--operation=PutProjectsProjectID \
#--operation=DeleteProjectsProjectID \
#--model=ProjectReq \
#--model=Project \
#--model=CVEWhitelist \
#--model=ProjectMetadata \
#--model=CVEWhitelistItem

# PostProjectsProjectIDRobots DeleteProjectsProjectIDRobotsRobotID GetProjectsProjectIDRobots
# RobotAccountCreate
