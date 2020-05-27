# Development

## Visual Studio Code DevContainer

For Easy development use [Visual Studio Code DevContainer](https://code.visualstudio.com/docs/remote/containers), you can find the basement from the Development Containers at [nolte/vscode-devcontainers](https://github.com/nolte/vscode-devcontainers).

### Precondition for Use

* Create you Github Personal Access Token under [github.com/settings/tokens](https://github.com/settings/tokens) with the following scopes:
    * `read:packages`

* Login to fetch the required dev containers

```sh
    docker login docker.pkg.github.com
```

* Grab you a Coffee and wait for 3 Minutes (This happens on the first time use)
* Click Terminal -> New Terminal and execute the following command:

```sh
    make
```

## Process

| *branch*                 | *description*                                                                    |
|--------------------------|----------------------------------------------------------------------------------|
| ```master```  | The ```master``` will only used for Presentation, and the Deployment process.    |
| ```tags/release/``` | All Created Releases must be start with a ```release/``` prefix at the tag name. |
| ```gh-pages``` |  |
| ```develop``` |  |
| ```feature/*``` |  |
| ```fix/*``` |  |


Please use the ```develop``` branch for new features and fixes.


### Releasing

The [Github Release](https://github.com/nolte/terraform-provider-harbor/releases) Assets will be automatic attach from the build job see ```.github/workflows/go.yml```.
![Go](https://github.com/nolte/terraform-provider-harbor/workflows/Go/badge.svg?branch=master)  
For a Easy Release process we use the GitHub Commandline Interface [cli.github.com](https://cli.github.com/manual/).  
Each Release will be start from the ```develop``` branch.

```sh
TBD
```

## Docs

**starting**
```bash
mkdocs serve
```
and open [127.0.0.1:8000](http://127.0.0.1:8000/)

## Development Shortcuts

**build and test in one command**

```sh
make compile \
    && make install \
    && bats scripts/test/bats/build
```

```
terraform import -var harbor_endpoint=${HARBOR_ENDPOINT} -var harbor_base_path='/api' harbor_project.main 24
```

## Links

* [writing-custom-providers](https://www.terraform.io/docs/extend/writing-custom-providers.html)
