# Development

## Visual Studio Code DevContainer

For Easy development use [Visual Studio Code DevContainer](https://code.visualstudio.com/docs/remote/containers), you can find the basement from the Development Containers at [nolte/vscode-devcontainers](https://github.com/nolte/vscode-devcontainers).

### Percondition for Use

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

The [Github Release](https://github.com/nolte/terraform-provider-harbor/releases) Assets will be automatical attatch from the build job see ```.github/workflows/go.yml```. 
![Go](https://github.com/nolte/terraform-provider-harbor/workflows/Go/badge.svg?branch=master)


## Docs

**starting**
```bash
mkdocs serve
```
and open [127.0.0.1:8000](http://127.0.0.1:8000/)