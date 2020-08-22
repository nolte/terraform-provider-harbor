//+build mage

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/sh"
	copy "github.com/otiai10/copy"
)

func GenerateDocumation() error {
	documentationTemplatesDir := "./provider_doc"
	outputDir := "../docs"
	mkdocsOut := "./mkdocs/gen"
	file, err := ioutil.TempDir("", "example")
	log.Printf("file %s", file)
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(file)
	err = filepath.Walk(documentationTemplatesDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			log.Print(path)
			if strings.HasSuffix(path, ".mdpp") {
				fmt.Println(path, info.Size())
				targetPath := strings.ReplaceAll(filepath.Join(file, path), ".mdpp", ".md")
				command := fmt.Sprintf("%s -o %s", path, targetPath)
				os.MkdirAll(filepath.Dir(targetPath), 0700)
				//log.Printf("Command %s", filepath.Dir(targetPath))
				err = sh.Run("markdown-pp", strings.Split(command, " ")...)
				if err != nil {
					log.Fatal(err)
				}
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	copyCommand := fmt.Sprintf("%s/provider_doc", file)
	err = copy.Copy(copyCommand, outputDir)
	err = copy.Copy(copyCommand, mkdocsOut)

	fmt.Printf("COPY: %s", copyCommand)
	return err
}
