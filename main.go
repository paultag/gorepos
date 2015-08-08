package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

// {{{ Package Struct
type Package struct {
	Repo     string
	Path     string
	Packages []string
	Url      string
}

func (p Package) GetRoutes() (ret []string) {
	ret = append(ret, p.Path)
	for _, pkg := range p.Packages {
		ret = append(ret, p.Path+"/"+pkg)
	}
	return
}

// }}}

// {{{ JSON Functions

func writePage(w http.ResponseWriter, data string, code int) error {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(code)
	w.Write([]byte(data))
	return nil
}

func writePageError(w http.ResponseWriter, message string, code int) error {
	return writePage(w, message, code)
}

func writePageGo(w http.ResponseWriter, path string, pkg Package, code int) error {
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
	return writePage(w, page.String(), code)
}

// }}}

// {{{ Config Functions

func loadConfig() (ret []Package, err error) {
	content, err := ioutil.ReadFile("packages.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &ret)
	return
}

func loadRoutes(packages []Package) (routes map[string]Package) {
	routes = map[string]Package{}

	for _, pkg := range packages {
		for _, route := range pkg.GetRoutes() {
			routes[route] = pkg
		}
	}
	return
}

// }}}

func main() {
	mux := http.NewServeMux()
	packages, err := loadConfig()
	if err != nil {
		fmt.Errorf("%s\n", err)
		return
	}
	routes := loadRoutes(packages)

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		route := req.URL.Path
		var err error
		if val, ok := routes[route]; ok {
			err = writePageGo(w, "pault.ag"+val.Path, val, 200)
		} else {
			err = writePageError(w, ":(", 404)
		}
		if err != nil {
			writePageError(w, fmt.Sprintf("error: %v", err), 400)
		}
	})
	http.ListenAndServe(":8000", mux)
}

// vim: foldmethod=marker
