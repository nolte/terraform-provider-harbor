# Harbor Provider

The ```terraform-provider-harbor``` is used to configure an instance of [Harbor](https://goharbor.io) in version `1.10.2`

This is original based on the Work from [BESTSELLER/terraform-harbor-provider](https://github.com/BESTSELLER/terraform-harbor-provider), but with some incompatible changes, like the access to the Harbor API.

![Classic CI/CD](https://github.com/nolte/terraform-provider-harbor/workflows/Classic%20CI/CD/badge.svg)
![Release Flow](https://github.com/nolte/terraform-provider-harbor/workflows/Release%20Flow/badge.svg)

## Project Status

**At the Moment this Project ist heavy under Construction, it is not recommendet for Production use, ~~or active Forking~~ !**

**Planed Branking Changes:**

- [ ] Rename provier attributes, like url etc.
- [x] Planed Git Rebase for remove the Ugly CI/CD Test Commit
  - [x] Finazilize the frist version of common ci workflow
  - [x] Finazilize the frist version of release workflow
- [x] use a ```develop``` branch as default
- [x] cleanup unused stuff from starting development
- [ ] Use First Stable version from the Devcontainer [docker.pkg.github.com/nolte/vscode-devcontainers/k8s-operator:latest](https://github.com/nolte/vscode-devcontainers) _(not exists at the moment)_

## Docs

The Documentation will be created with [mkdocs](https://www.mkdocs.org/) and generated to [nolte.github.io/terraform-provider-harbor](https://nolte.github.io/terraform-provider-harbor/) from the latest Release like, ![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/nolte/terraform-provider-harbor).

## Building

As CI/CD tool we use the Github Workflow Feature.

## Visual Studio Code DevContainer

For Easy development use [Visual Studio Code DevContainer](https://code.visualstudio.com/docs/remote/containers), you can find the basement from the Development Containers at [nolte/vscode-devcontainers](https://github.com/nolte/vscode-devcontainers).

1. Create you Github Personal Access Token under <https://github.com/settings/tokens> with the following scopes:
   1. `read:packages`

2. Login to fetch the required dev containers

    ```sh
    docker login docker.pkg.github.com
    ```

3. Grab you a Coffee and wait for 3 Minutes (This happens on the first time use)

4. Click Terminal -> New Terminal and execute the following command:

```sh
# using the Makefile
make
```

### Precondition Tools

For full building and testing you need the following tools on our machine.

#### Required For Building

- go
- [go-swagger/go-swagger](https://github.com/go-swagger/go-swagger)
- [WindomZ/swagger-merger](https://github.com/WindomZ/swagger-merger)

#### Required For Testing

- kind
- Terraform
- bats

## Supported Versions

Tested with Harbor v2.1.0 and v1.10.2.

## Tests

For the End To End Tests we use a local [kind](https://kind.sigs.k8s.io) _(KubernetesInDocker)_ Cluster.
