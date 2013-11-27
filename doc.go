package main

import (
	"doc"
	"gnd.la/mux"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type packageGroup struct {
	Title    string
	Packages []*doc.Package
}

func DocListHandler(ctx *mux.Context) {
	pkgs, err := doc.ImportPackages(filepath.Join(doc.SourceDir, "gnd.la"))
	if err != nil {
		panic(err)
	}
	groups := []packageGroup{
		{Title: "Gondola Packages", Packages: pkgs},
	}
	infos, _ := ioutil.ReadDir(doc.SourceDir)
	var opkgs []*doc.Package
	for _, v := range infos {
		if v.IsDir() && v.Name() != "gnd.la" {
			p, _ := doc.ImportPackages(filepath.Join(doc.SourceDir, v.Name()))
			opkgs = append(opkgs, p...)
		}
	}
	if len(opkgs) > 0 {
		groups = append(groups, packageGroup{
			Title:    "Other useful Rainy Cape's packages",
			Packages: opkgs,
		})
	}
	title := "Package Index"
	data := map[string]interface{}{
		"HeaderTitle": title,
		"Subtitle":    title,
		"Section":     "docs",
		"Groups":      groups,
	}
	ctx.MustExecute("pkgs.html", data)
}

func DocHandler(ctx *mux.Context) {
	path := ctx.IndexValue(0)
	pkg, err := doc.ImportPackage(path)
	if err != nil {
		panic(err)
	}
	title := "Package " + pkg.ImportPath()
	breadcrumbs := []*Breadcrumb{
		{Title: "Index", Href: ctx.MustReverse("doc-list")},
	}
	for ii := 0; ii < len(path); {
		var end int
		slash := strings.IndexByte(path[ii:], '/')
		if slash < 0 {
			end = len(path)
		} else {
			end = ii + slash
		}
		breadcrumbs = append(breadcrumbs, &Breadcrumb{
			Title: path[ii:end],
			Href:  ctx.MustReverse("doc", path[:end]),
		})
		ii = end + 1
	}
	data := map[string]interface{}{
		"HeaderTitle": "Package " + pkg.Name(),
		"Subtitle":    title,
		"Section":     "docs",
		"Breadcrumbs": breadcrumbs,
		"Package":     pkg,
	}
	ctx.MustExecute("package.html", data)
}
