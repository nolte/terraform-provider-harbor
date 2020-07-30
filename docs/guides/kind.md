# E2E Tests With Kind

For Quick and Easy Local Development it is Recommended to use a Vanilla [Harbor]() installation. The Local Test environment will be used the [Helm Chart](https://helm.sh/) from [goharbor/harbor-helm](https://github.com/goharbor/harbor-helm/releases).  

All relevant `Makefile` goals are prefixed with ```e2s_*```.

| Goal                    | Description                                                                                  |
|-------------------------|----------------------------------------------------------------------------------------------|
| `e2e_prepare`           | Configure the Kind based Env, and install the latest version from the Harbor Chart.          |
| `e2e_prepare_harbor_v1` | Install a Harbor v1 to the Kind Cluster                                                      |
| `e2e_prepare_harbor_v2` | Install a Harbor v2 to the Kind Cluster                                                      |
| `e2e_test_v1`           | Starting the go based tests again a Harbor V1 Deployment                                     |
| `e2e_test_v2`           | Starting the go based tests again a Harbor V2 Deployment                                     |
| `e2e_test_classic`      | Test the Terraform Scripts from `examples` with [terratest](https://terratest.gruntwork.io/) |
| `e2e_clean_harbor`      | Delete the Harbor helm chart from kind cluster                                               |
| `e2e_clean_cluster`     | Remove the current kind cluster.                                                             |

The same flavore of tests will be integrated into the [CI/CD](/guides/development/#cicd) Process.


**Befor you execute `e2e_prepare` ensure that the `INGRESS_DOMAIN` to your local host** (By default we use the Docker Bridge Network Interface like `172-17-0-1.sslip.io`.)

*example without change any files*

```bash
./scripts/tst-00-prepare-kind.sh "2.1.0"
./scripts/tst-01-prepare-harbor.sh "10-42-0-100.sslip.io" "1.2.0"
```

As wildcard DNS we use a service like [sslip.io](https://sslip.io/).  
The scripts will be starting the [Kind](https://kind.sigs.k8s.io) Cluster with [Ingress Support](https://kind.sigs.k8s.io/docs/user/ingress/), as Ingress Controller we use [Nginx](https://kind.sigs.k8s.io/docs/user/ingress/#ingress-nginx) .

The [Kind](https://kind.sigs.k8s.io) installation is part of the used [Visual Studio Code Devcontainer](/guides/development/#visual-studio-code-devcontainer)

## Go Based Tests

For quick response and tests near the Code we use [Terrafomr Go Acceptance Tests](https://www.terraform.io/docs/extend/testing/acceptance-tests/index.html).

```bash
make e2e_test
```

Tests will be matchs by the file name prefix `*_test.go`.

## Terratest File Based Tests

The Classic Terraform File based tests are located at `examples/**`, and executed with [terratest](https://terratest.gruntwork.io/). For execution you need, a full runnable Terraform Enviroment with the current version from the harbor terraform provider.

```bash
# compile the provider
make compile

# copy provider to local ~/.terraform.d/plugins/linux_amd64 folder
make install

# start the tests
make e2e_test_classic
```
