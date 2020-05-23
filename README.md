# Harbor Provider
The ```terraform-provider-harbor``` is used to configure an instance of [Harbor](https://goharbor.io).

This Frok is original based on the Work from [BESTSELLER/terraform-harbor-provider](https://github.com/BESTSELLER/terraform-harbor-provider), but with some incompatible changes, like the access to the Harbor API.

![Go](https://github.com/nolte/terraform-provider-harbor/workflows/Go/badge.svg)

**Project Status**

**At the Moment this Project ist heavy under Construction, it is not recommendet for Production use, or active Forking!**

**Planed Branking Changes:**
- [ ] Rename provier attributes, like url etc.
- [ ] Planed Git Rebase for remove the Ugly CI/CD Test Commit
    - [ ] Finazilize the frist version of common ci workflow
    - [ ] Finazilize the frist version of release workflow
- [ ] use a ```develop``` branch as default
- [ ] cleanup unused stuff from starting development

## Docs

The Documentation will be created with [mkdocs](https://www.mkdocs.org/) and generated to [nolte.github.io/terraform-provider-harbor](https://nolte.github.io/terraform-provider-harbor/) from the latest Release like, ![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/nolte/terraform-provider-harbor). 

## Building

As CI/CD tool we use the Github Workflow Feature.

## Visual Studio Code DevContainer

For Easy development use [Visual Studio Code DevContainer](https://code.visualstudio.com/docs/remote/containers), you can find the basement from the Development Containers at [nolte/vscode-devcontainers](https://github.com/nolte/vscode-devcontainers).

1. Create you Github Personal Access Token under https://github.com/settings/tokens with the following scopes:
   1. `read:packages`

2. Login to fetch the required dev containers

```sh
docker login docker.pkg.github.com
```

4. Grab you a Coffee and wait for 3 Minutes (This happens on the first time use)

3. Click Terminal -> New Terminal and execute the following command:

```sh
# using the Makefile
make
```

### Precondition Tools

For full building and testing you need the following tools on our machine.

**Required For Building**

* go
* swagger
* swagger-merger

**Required For Testing**
* kind
* Terraform
* bats

## Supported Versions

Tested with Harbor v2.1.0 and v1.10.2.

## Tests

For the End To End Tests we use a local [kind](https://kind.sigs.k8s.io) _(KubernetesInDocker)_ Cluster.

