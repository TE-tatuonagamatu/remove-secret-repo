package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Masterminds/glide/cfg"
	"github.com/Masterminds/glide/importer"
	yaml "gopkg.in/yaml.v2"
)

func getPackageName() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	gopath := os.Getenv("GOPATH")
	if strings.HasPrefix(cwd, gopath) {
		return strings.Replace(strings.Replace(cwd, gopath, "", 1), "/src/", "", 1)
	}

	return cwd
}

func filterDependencies(dependencies cfg.Dependencies) cfg.Dependencies {
	var dst cfg.Dependencies
	for _, d := range dependencies {
		if strings.HasSuffix(d.Name, "-secrets") {
			fmt.Printf("Remove %s\n", d.Name)
			continue
		}
		dst = append(dst, d)
	}
	return dst
}

func main() {
	exist, dependencies, err := importer.Import(".")
	if err != nil {
		panic("Failed to import dependency file")
	}

	if !exist {
		panic("No dependency config file found")
	}

	b, err := yaml.Marshal(filterDependencies(dependencies))
	if err != nil {
		panic("Failed to marshal glide.yaml")
	}

	f, err := ioutil.TempFile(".", "glide")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	hdrTemplate := "package: %s\nimport:\n"
	_, err = f.WriteString(fmt.Sprintf(hdrTemplate, getPackageName()))
	if err != nil {
		panic(err)
	}

	_, err = f.Write(b)
	if err != nil {
		panic(err)
	}

	err = os.Rename(f.Name(), "glide.yaml")
	if err != nil {
		panic(err)
	}
}
