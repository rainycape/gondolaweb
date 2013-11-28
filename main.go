package main

import (
	"doc"
	"gnd.la/admin"
	_ "gnd.la/bootstrap"
	"gnd.la/config"
	"gnd.la/mux"
	_ "gnd.la/orm/driver/postgres"
	"gnd.la/util"
)

var (
	STATIC_FILES_PATH = util.RelativePath("static")
)

var Config struct {
	config.Config
}

type Breadcrumb struct {
	Title string
	Href  string
}

func main() {
	doc.SourceDir = util.RelativePath("src")
	doc.Context.GOPATH = util.RelativePath(".")
	doc.SourceHandlerPrefix = "/doc/src/"
	doc.DocHandlerPrefix = "/doc/pkg/"
	config.MustParse(&Config)
	StartUpdatingRepos()
	m := mux.New()
	m.HandleFunc("^/$", mux.TemplateHandler("main.html", util.M{"Section": "home"}))
	m.HandleNamedFunc("^/tutorials/$", "tutorials", nil)
	m.HandleNamedFunc("^/tutorial/([\\w\\-]+)/$", "tutorial", nil)
	m.HandleFunc("^/article/$", mux.TemplateHandler("article.html", util.M{"Section": "docs"}))
	m.HandleNamedFunc("^/doc/src/(.+)", "source", SourceHandler)
	m.HandleNamedFunc("^/doc/pkg/$", "doc-list", DocListHandler)
	m.HandleNamedFunc("^/doc/pkg/std/$", "std-doc-list", DocStdListHandler)
	m.HandleNamedFunc("^/doc/pkg/(.+)", "doc", DocHandler)
	m.SetTrustXHeaders(true)
	m.HandleAssets("/static/", STATIC_FILES_PATH)
	if !admin.Perform(m) {
		m.MustListenAndServe()
	}
}
