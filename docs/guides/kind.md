# E2E Tests With Kind

For Quick and Easy Local Development it is Recommended to use a Vanilla Harbor installation.
All relevant make goals are prefixed with ```e2s_*```.

Makefile Goals

| Goal          | Description                                                                         |
|---------------|-------------------------------------------------------------------------------------|
| `e2e_prepare` | Configure the Kind based Env, and install the latest version from the Harbor Chart. |
| `e2e_test`    | Starting the go based tests again a Harbor Deployment                               |
| `e2e_test_classic` | Remove the current kind cluster.                                                    |
| `e2e_cleanup` | Remove the current kind cluster.                                                    |



**Befor you execute `e2e_prepare` ensure that the `INGRESS_DOMAIN` default ar configured to your local host**

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

## Terraform File Based Tests

The Classic Terraform File based tests are located at `scripts/test/tf-acception-test`.

```bash
# compile the provider
make compile 

# copy provider to local ~/.terraform.d/plugins/linux_amd64 folder
make install

# start the tests
make e2e_test_classic
```
