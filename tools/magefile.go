//+build mage

package main

import (
	"context"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
	flags "github.com/jessevdk/go-flags"
	"io/ioutil"
	"os"
	"path"
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

    args = append(args,"--with-flatten=remove-unused")



	return clt.Execute(args)
}
