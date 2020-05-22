# Harbor Provider
The ```terraform-provider-harbor``` is used to configure an instance of [Harbor](https://goharbor.io).

This Frok is original based on the Work from [BESTSELLER/terraform-harbor-provider](https://github.com/BESTSELLER/terraform-harbor-provider), but with some incompatible changes, like the access to the Harbor API.

![Go](https://github.com/nolte/terraform-provider-harbor/workflows/Go/badge.svg)

## Docs

The Documentation will be created with [mkdocs](https://www.mkdocs.org/) and generated to [nolte.github.io/terraform-provider-harbor](https://nolte.github.io/terraform-provider-harbor/).

**starting**
```bash
mkdocs serve
```
and open [127.0.0.1:8001](http://127.0.0.1:8001/)

## Building

For Easy development use [Visual Studio Code DevContainer](https://code.visualstudio.com/docs/remote/containers), you can find the basement from the Development Containers at [nolte/vscode-devcontainers](https://github.com/nolte/vscode-devcontainers).

```bash
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

