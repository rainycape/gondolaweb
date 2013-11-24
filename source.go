package main

import (
	"gnd.la/mux"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	highlighters = map[string]string{
		".c":    "c",
		".cpp":  "cpp",
		".css":  "css",
		".cxx":  "cpp",
		".go":   "go",
		".h":    "c",
		".hpp":  "cpp",
		".hxx":  "cpp",
		".js":   "js",
		".sh":   "sh",
		".bash": "sh",
	}
)

func SourceHandler(ctx *mux.Context) {
	rel := ctx.IndexValue(0)
	var path string
	if strings.HasPrefix(rel, "go") {
		path = filepath.Join(buildContext().GOROOT, "src", filepath.FromSlash(rel[2:]))
	} else {
		path = filepath.Join(srcDir, filepath.FromSlash(rel))
	}
	info, err := os.Stat(path)
	if err != nil {
		ctx.NotFound("File not found")
		return
	}
	var breadcrumbs []*Breadcrumb
	for ii := 0; ii < len(rel); {
		var end int
		slash := strings.IndexByte(rel[ii:], '/')
		if slash < 0 {
			end = len(rel)
		} else {
			end = ii + slash
		}
		breadcrumbs = append(breadcrumbs, &Breadcrumb{
			Title: rel[ii:end],
			Href:  ctx.MustReverse("source", rel[:end]),
		})
		ii = end + 1
	}
	var template string
	var title string
	var files []string
	var code string
	if info.IsDir() {
		if rel != "" && rel[len(rel)-1] != '/' {
			ctx.MustRedirectReverse(true, "source", rel+"/")
			return
		}
		contents, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, v := range contents {
			if n := v.Name(); len(n) > 0 && n[0] != '.' {
				files = append(files, n)
			}
		}
		title = "Directory " + filepath.Base(rel)
		template = "dir.html"
	} else {
		contents, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		title = "File " + filepath.Base(rel)
		code = string(contents)
		template = "source.html"
	}
	data := map[string]interface{}{
		"Sections":    "docs",
		"Subtitle":    rel,
		"HeaderTitle": title,
		"Breadcrumbs": breadcrumbs,
		"Files":       files,
		"Code":        code,
		"Highlighter": highlighters[filepath.Ext(rel)],
	}
	ctx.MustExecute(template, data)
}
