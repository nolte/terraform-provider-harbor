//+build mage

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
	_ "github.com/golangci/golangci-lint/pkg/commands"
	"github.com/goreleaser/goreleaser/cmd"
	flags "github.com/jessevdk/go-flags"
	"github.com/magefile/mage/mg"
	"github.com/nolte/plumbing/cmd/golang"

	// mage:import
	_ "github.com/nolte/plumbing/cmd/kind"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func GenerateHarborGoClient(ctx context.Context) error {
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

func Lint(ctx context.Context) {
	ctx = context.WithValue(ctx, "basedir", "../")
	mg.CtxDeps(ctx, golang.Golang.Lint)
}

func Fmt(ctx context.Context) {
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

func GoRelease() {
	os.Chdir("../")
	defer os.Chdir("./tools")
	args := []string{"build", "--rm-dist", "--snapshot"}
	cmd.Execute(
		buildVersion(version, commit, date, builtBy),
		os.Exit,
		args,
	)
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
