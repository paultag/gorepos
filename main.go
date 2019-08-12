package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path"
)

type Package struct {
	Repo     string
	Path     string
	Packages []string
	Url      string
}

func writePage(w io.Writer, path string, pkg Package) error {
	t, err := template.ParseFiles("template.html")
	if err != nil {
		return err
	}
	pkg.Path = path
	var page bytes.Buffer
	err = t.Execute(&page, pkg)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, &page)
	return err
}

func loadConfig() (ret []Package, err error) {
	content, err := ioutil.ReadFile("packages.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &ret)
	return
}

func main() {
	packages, err := loadConfig()
	if err != nil {
		fmt.Errorf("%s\n", err)
		return
	}

	namespace := "pault.ag"

	for _, pkg := range packages {
		packages := pkg.Packages
		packages = append(packages, "")
		for _, subpackage := range packages {
			pkgpath := path.Join(namespace, pkg.Path)

			root := path.Join(".", pkg.Path, subpackage)

			if err := os.MkdirAll(root, 0755); err != nil {
				fmt.Errorf("%s\n", err)
				return
			}

			fd, err := os.Create(path.Join(root, "index.html"))
			if err != nil {
				fmt.Errorf("%s\n", err)
				return
			}
			defer fd.Close()
			writePage(fd, pkgpath, pkg)
			fd.Close()
		}
	}
}

// vim: foldmethod=marker
