# Harbor Provider

The ```terraform-provider-harbor``` is used to configure an instance of [Harbor](https://goharbor.io).

This is original based on the Work from [BESTSELLER/terraform-harbor-provider](https://github.com/BESTSELLER/terraform-harbor-provider), but with some incompatible changes, like the access to the Harbor API.

[![Classic CI/CD](https://github.com/nolte/terraform-provider-harbor/workflows/Classic%20CI/CD/badge.svg)](https://github.com/nolte/terraform-provider-harbor/actions?query=workflow%3A%22Classic+CI%2FCD%22)
[![Release Flow](https://github.com/nolte/terraform-provider-harbor/workflows/Release%20Flow/badge.svg)](https://github.com/nolte/terraform-provider-harbor/actions?query=workflow%3A%22Release+Flow%22)



## Project Status

**At the Moment this Project is heavy under Construction, it is not recommended for Production use, ~~or active Forking~~ !**

**Planed Braking Changes:**

- [x] [Rename provider](https://github.com/nolte/terraform-provider-harbor/issues/3) attributes, like url etc.
- [ ] [Refectore Config Auth](https://github.com/nolte/terraform-provider-harbor/issues/10) attributes, like ldap etc.
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

After starting the VisualStudio Code DevContainer, you can access the Documentation at [localhost:8000](http://localhost:8000).

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

Tested with Harbor v1.10.2, v2.0.0 and v2.1.0.
