//+build mage

package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	semver "github.com/blang/semver/v4"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
	_ "github.com/golangci/golangci-lint/pkg/commands"
	"github.com/goreleaser/goreleaser/cmd"
	flags "github.com/jessevdk/go-flags"
	"github.com/magefile/mage/mg"
	"github.com/nolte/plumbing/cmd/golang"

	// mage:import
	_ "github.com/nolte/plumbing/cmd/kind"

	plumbing "github.com/nolte/plumbing/pkg"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func newHarborHelmDeployment(version string) (deployment plumbing.HelmDeployment, selectors map[string]string) {
	targetNamespace := "harbor"
	ingressDnsName := os.Getenv("INGRESS_DOMAIN")
	if ingressDnsName == "" {
		ingressDnsName = "172-17-0-1.sslip.io"
	}

	chart := plumbing.HelmDeployment{
		ExtraValues: map[string]string{
			"expose.ingress.hosts.core":   fmt.Sprintf("harbor.%s", ingressDnsName),
			"expose.ingress.hosts.notary": fmt.Sprintf("notary.%s", ingressDnsName),
			"externalURL":                 fmt.Sprintf("https://harbor.%s", ingressDnsName),
		},
		ReleaseName: "tf-harbor-test",
		Namespace:   targetNamespace,
		Chart: plumbing.HelmChart{
			Name:    "harbor",
			Version: version,
			Repository: plumbing.HelmRepository{
				Name: "harbor",
				URL:  "https://helm.goharbor.io",
			},
		},
	}
	labels := map[string]string{
		"app": "harbor",
	}
	return chart, labels
}

// TestArtefacts configure the Test Artefacts.
type TestArtefacts mg.Namespace

// Deploy Harbor Helm Chart to Cluster.
func (TestArtefacts) Deploy(ctx context.Context) {
	version := os.Getenv("HARBOR_HELM_CHART_VERSION")
	chart, labels := newHarborHelmDeployment(version)
	plumbing.ApplyHelmChart(chart, labels)
}

// DeployV1 Harbor Helm Chart to Cluster.
func (TestArtefacts) DeployV1(ctx context.Context) {
	chart, labels := newHarborHelmDeployment("1.3.2")
	plumbing.ApplyHelmChart(chart, labels)
}

// DeployV2 Harbor Helm Chart to Cluster.
func (TestArtefacts) DeployV2(ctx context.Context) {
	chart, labels := newHarborHelmDeployment("1.4.2")
	plumbing.ApplyHelmChart(chart, labels)
}

// Delete Harbor Helm Chart from Cluster.
func (TestArtefacts) Delete(ctx context.Context) {
	chart, _ := newHarborHelmDeployment("1.4.2")
	chart.Delete()
}

// Build configure the Build Targets.
type Build mg.Namespace

func (Build) GenerateHarborGoClient(ctx context.Context) error {
	originalDat, err := ioutil.ReadFile("../scripts/swagger-specs/v2-swagger-original.json")
	check(err)
	patchDat, err := ioutil.ReadFile("../scripts/swagger-specs/patch.1.json")
	check(err)
	patch, err := jsonpatch.DecodePatch(patchDat)
	if err != nil {
		panic(err)
	}

	modified, err := patch.Apply(originalDat)
	if err != nil {
		panic(err)
	}

	configPath := path.Join(os.TempDir(), "patched-swagger.json")
	err = ioutil.WriteFile(configPath, modified, 0o600)
	check(err)
	defer os.Remove(configPath)
	return generateGoSourcesFromSwaggerSpec(configPath)
}

func generateGoSourcesFromSwaggerSpec(path string) error {

	generatedPath := "../gen/harborctl"
	os.RemoveAll(generatedPath)
	err := os.MkdirAll(generatedPath, 0777)
	check(err)

	clt := generate.Client{}
	clt.Shared.Spec = flags.Filename(path)
	clt.Shared.Target = flags.Filename(generatedPath)
	clt.Models.ModelPackage = "models"
	clt.Name = "harbor"
	clt.ClientPackage = "client"

	var args []string

	args = append(args, "--with-flatten=remove-unused")

	return clt.Execute(args)
}

func (Build) Lint(ctx context.Context) {
	ctx = context.WithValue(ctx, "basedir", "../")
	mg.CtxDeps(ctx, golang.Golang.Lint)
}

func (Build) Fmt(ctx context.Context) {
	ctx = context.WithValue(ctx, "basedir", "../")
	mg.CtxDeps(ctx, golang.Golang.Fmt)
}

// nolint: gochecknoglobals
var (
	version = "dev"
	commit  = ""
	date    = ""
	builtBy = ""
)

func (Build) GoRelease() {
	os.Chdir("../")
	defer os.Chdir("./tools")
	args := []string{"build", "--rm-dist", "--snapshot"}
	cmd.Execute(
		buildVersion(version, commit, date, builtBy),
		os.Exit,
		args,
	)
}

func Copy(src, dst string) error {
	in, err := os.Open(src)
	check(err)
	defer in.Close()
	destDir, err := filepath.Abs(filepath.Dir(dst))
	check(err)

	err = os.MkdirAll(destDir, 0700)
	check(err)

	out, err := os.Create(dst)
	check(err)
	defer out.Close()

	_, err = io.Copy(out, in)
	check(err)

	err = os.Chmod(dst, 0744)
	check(err)

	return out.Close()
}

type TerraformInstallation struct {
	pluginsDir       string
	providerFileName string
}

func (tf TerraformInstallation) InstallPlugin(providerPath string) error {
	dest := filepath.Join(tf.pluginsDir, tf.providerFileName)
	log.Printf("Copy privider to %s", dest)
	return Copy(providerPath, dest)
}

func (tf TerraformInstallation) CleanPlugins() error {
	os.RemoveAll(tf.pluginsDir)
	return os.MkdirAll(tf.pluginsDir, 0700)
}

func tf13(home string) TerraformInstallation {
	return TerraformInstallation{
		pluginsDir:       filepath.Join(home, ".local/share/terraform/plugins/registry.terraform.io/nolte/harbor/0.0.1/linux_amd64"),
		providerFileName: "terraform-provider-harbor",
	}
}

func tf12(home string) TerraformInstallation {
	return TerraformInstallation{
		pluginsDir:       filepath.Join(home, ".terraform.d/plugins/linux_amd64"),
		providerFileName: "terraform-provider-harbor_v0.0.1",
	}
}

func terraformCustomDataDir() string {
	home, err := os.UserHomeDir()
	check(err)
	return filepath.Join(home, ".local/share/terraform/plugins")
}

func (Build) TerraformInstallProvider() {

	distPath := "../dist/terraform-provider-harbor_linux_amd64"
	files, err := ioutil.ReadDir(distPath)
	check(err)
	for _, f := range files {
		localFile := filepath.Join(distPath, f.Name())

		tfPlugins := terraformPluginDir()
		tfPlugins.CleanPlugins()
		tfPlugins.InstallPlugin(localFile)

	}
}

func terraformPluginDir() TerraformInstallation {
	versionTxt := os.Getenv("TF_VERSION")
	if versionTxt == "" {
		versionTxt = "0.13.0"
	}
	version, err := semver.Make(versionTxt)
	home, err := os.UserHomeDir()
	check(err)
	v13, err := semver.Make("0.13.0")
	check(err)
	if v13.Compare(version) == 0 {
		return tf13(home)
	} else {
		return tf12(home)
	}
}

func buildVersion(version, commit, date, builtBy string) string {
	var result = version
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}
	if builtBy != "" {
		result = fmt.Sprintf("%s\nbuilt by: %s", result, builtBy)
	}
	return result
}
