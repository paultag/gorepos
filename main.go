package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// {{{ Package Struct
type Package struct {
	Repo     string
	Path     string
	Packages []string
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

func writePageGo(w http.ResponseWriter, path string, pkg *Package, code int) error {
	return writePage(w, fmt.Sprintf(`<!DOCTYPE html><html>
	<head><meta charset="utf-8"><title>%s</title><meta name="go-import" content="%s %s %s"></head>
    <body>%s</body>
</html>`, "Package!", path, "git", pkg.Repo, "Package!"), code)
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

func loadRoutes(packages []Package) (routes map[string]*Package) {
	routes = map[string]*Package{}

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
		if val, ok := routes[route]; ok {
			writePageGo(w, req.Host+route, val, 200)
		} else {
			writePageError(w, ":(", 401)
		}
	})
	http.ListenAndServe(":8000", mux)
}

// vim: foldmethod=marker
