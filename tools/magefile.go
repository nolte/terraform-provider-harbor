//+build mage

package main

import (
	"context"
    "io/ioutil"
    "os"
	"path"
    jsonpatch "github.com/evanphx/json-patch"
    "github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
    flags "github.com/jessevdk/go-flags"

)
func check(e error) {
    if e != nil {
        panic(e)
    }
}

func TestCommand(ctx context.Context) error {
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

    clt := generate.Client{

    }
	clt.Shared.Spec = flags.Filename(path)
    clt.Shared.Target = flags.Filename("../gen/harborctl")
    //clt.Shared.Name = "harbor"

    var args []string

    return clt.Execute(args)
}
