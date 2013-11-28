package main

import (
	"bytes"
	"doc"
	"fmt"
	"gnd.la/html"
	"gnd.la/log"
	"gnd.la/mux"
	"html/template"
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
		".md":   "markdown",
		".sh":   "sh",
		".bash": "sh",
	}
)

func SourceHandler(ctx *mux.Context) {
	rel := ctx.IndexValue(0)
	var path string
	if strings.IndexByte(rel, '.') < strings.IndexByte(rel, '/') {
		// Non std package
		path = filepath.Join(doc.SourceDir, filepath.FromSlash(rel))
	} else {
		// Std pckage
		path = filepath.Join(doc.Context.GOROOT, "src", "pkg", filepath.FromSlash(rel))
	}
	log.Debugf("Loading source from %s", path)
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
	var tmpl string
	var title string
	var files []string
	var code template.HTML
	var lines []int
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
		tmpl = "dir.html"
	} else {
		contents, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		title = "File " + filepath.Base(rel)
		var buf bytes.Buffer
		buf.WriteString("<span id=\"line-1\">")
		last := 0
		line := 1
		for ii, v := range contents {
			if v == '\n' {
				buf.WriteString(html.Escape(string(contents[last:ii])))
				lines = append(lines, line)
				last = ii
				line++
				buf.WriteString(fmt.Sprintf("</span><span id=\"line-%d\">", line))
			}
		}
		buf.Write(contents[last:])
		buf.WriteString("</span>")
		code = template.HTML(buf.String())
		tmpl = "source.html"
	}
	data := map[string]interface{}{
		"Sections":    "docs",
		"Subtitle":    rel,
		"HeaderTitle": title,
		"Breadcrumbs": breadcrumbs,
		"Files":       files,
		"Code":        code,
		"Lines":       lines,
		"Highlighter": highlighters[filepath.Ext(rel)],
	}
	ctx.MustExecute(tmpl, data)
}
