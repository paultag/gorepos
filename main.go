package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
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

func loadConfig(namespace string) (ret []Package, err error) {
	content, err := ioutil.ReadFile(fmt.Sprintf("%s.json", namespace))
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &ret)
	return
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("main.go namespace fs-root-dir\n")
	}

	namespace := os.Args[1]
	root := os.Args[2]

	packages, err := loadConfig(namespace)
	if err != nil {
		log.Fatalf("%s\n", err)
		return
	}

	for _, pkg := range packages {
		packages := pkg.Packages
		packages = append(packages, "")
		for _, subpackage := range packages {
			pkgpath := path.Join(namespace, pkg.Path)

			root := path.Join(root, pkg.Path, subpackage)

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
