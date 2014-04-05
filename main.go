package main

import (
	"gnd.la/app"
	"gnd.la/apps/articles"
	"gnd.la/apps/docs"
	_ "gnd.la/bootstrap"
	"gnd.la/config"
	_ "gnd.la/encoding/codec/msgpack"
	"gnd.la/internal/project"
	_ "gnd.la/orm/driver/postgres"
	_ "gnd.la/template/markdown"
	"gnd.la/util/pathutil"
	//	"time"
)

var (
	STATIC_FILES_PATH = pathutil.Relative("static")
	a                 *app.App
)

var Config struct {
	config.Config
}

func init() {
	config.MustParse(&Config)
	a = app.New()
	a.HandleAssets("/static/", STATIC_FILES_PATH)
	a.AddTemplateVars(map[string]interface{}{
		"Repo": "ssh://abra.rm-fr.net/home/fiam/git/gondola.git",
	})
	a.SetTrustXHeaders(true)

	// gnd.la handler, used by go get, etc...
	a.HandleOptions("/", gndlaHandler, &app.Options{Host: "gnd.la"})

	// Site handlers
	a.Handle("^/$", app.TemplateHandler("main.html", map[string]interface{}{"Section": "home"}))

	// docs app
	a.Include("/doc/", docs.App, "docs-base.html")
	docs.Groups = []*docs.Group{
		{"Gondola Packages", []string{"gnd.la/"}},
	}
	//docs.StartUpdatingPackages(time.Minute * 30)

	// articles app, clone it to load it twice: once
	// for the articles and once for the tutorials
	tutorialsApp := articles.App.Clone()
	tutorialsApp.SetName("Tutorials")
	a.Include("/tutorials/", tutorialsApp, "articles-base.html")
	articles.MustLoad(tutorialsApp, pathutil.Relative("tutorials"))

	articlesApp := articles.App.Clone()
	a.Include("/articles/", articlesApp, "articles-base.html")
	articles.MustLoad(articlesApp, pathutil.Relative("articles"))

	// API
	a.Handle("^/api/v1/templates$", app.JSONHandler(templateListHandler))
	a.HandleNamed("^/api/v1/template/download/([\\w\\-_]+)\\-v(\\d+)\\.tar\\.gz$", templateDownloadHandler, templateDownloadHandlerName)

	// Load project templates
	var err error
	templates, err = project.LoadTemplates(pathutil.Relative("templates"))
	if err != nil {
		panic(err)
	}
	for _, v := range templates {
		v.URL = a.MustReverse(templateDownloadHandlerName, v.Name, v.Version)
	}
}

func main() {
	a.MustListenAndServe()
}
