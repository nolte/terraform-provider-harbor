#!/usr/bin/env bash
# Bash3 Boilerplate. Copyright (c) 2014, kvz.io

set -o errexit
set -o pipefail
set -o nounset
# set -o xtrace

# Set magic variables for current file & dir
__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__file="${__dir}/$(basename "${BASH_SOURCE[0]}")"
__base="$(basename ${__file} .sh)"
__root="$(cd "$(dirname "${__dir}")" && pwd)" # <-- change this as it depends on your app

arg1="${1:-}"
projectBase=$__root

rm -rf ${projectBase}/gen/harborctl
rm ${projectBase}/gen/merged.json || true

swagger-merger \
    -o ${projectBase}/gen/merged.json \
    -i ${projectBase}/scripts/swagger-specs/v1-swagger-extra-fields.json \
    -i ${projectBase}/scripts/swagger-specs/v2-swagger-original.json 

mkdir -p ${projectBase}/gen/harborctl

swagger generate client \
    -f ${projectBase}/gen/merged.json \
    --name=harbor \
    --target=${projectBase}/gen/harborctl \
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
   