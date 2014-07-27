package main

import (
	"time"

	_ "gnd.la/admin"
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
	STATIC_FILES_PATH = pathutil.Relative("assets")
	App               *app.App
)

func init() {
	config.MustParse()
	App = app.New()
	App.HandleAssets("/static/", STATIC_FILES_PATH)
	App.AddTemplateVars(map[string]interface{}{
		"Repo": "ssh://abra.rm-fr.net/home/fiam/git/gondola.git",
	})
	App.SetTrustXHeaders(true)

	// gnd.la handler, used by go get, etc...
	App.HandleOptions("/", gndlaHandler, &app.Options{Host: "gnd.la"})

	// Site handlers
	App.Handle("^/$", app.TemplateHandler("main.html", map[string]interface{}{"Section": "home"}))

	// docs app
	docs.DefaultContext.GOROOT = goRoot
	docs.DefaultContext.GOPATH = goPath
	App.Include("/doc/", docs.App, "docs-base.html")
	docs.Groups = []*docs.Group{
		{"Gondola Packages", []string{"gnd.la/"}},
	}
	// Wait 10 seconds so the app starts and go get
	// can retrieve gnd.la, since this same app is
	// serving that content.
	time.AfterFunc(10*time.Second, func() {
		docs.StartUpdatingPackages(time.Minute * 10)
	})

	// articles app, clone it to load it twice: once
	// for the articles and once for the tutorials
	tutorialsApp := articles.App.Clone()
	if _, err := articles.LoadDir(tutorialsApp, pathutil.Relative("tutorials")); err != nil {
		panic(err)
	}
	tutorialsApp.SetName("Tutorials")
	App.Include("/tutorials/", tutorialsApp, "articles-base.html")

	articlesApp := articles.App.Clone()
	if _, err := articles.LoadDir(articlesApp, pathutil.Relative("articles")); err != nil {
		panic(err)
	}
	App.Include("/articles/", articlesApp, "articles-base.html")

	// API
	App.Handle("^/api/v1/templates$", app.JSONHandler(templateListHandler))
	App.HandleNamed("^/api/v1/template/download/([\\w\\-_]+)\\-v(\\d+)\\.tar\\.gz$", templateDownloadHandler, templateDownloadHandlerName)

	// Load project templates
	var err error
	templates, err = project.LoadTemplates(pathutil.Relative("templates"))
	if err != nil {
		panic(err)
	}
	for _, v := range templates {
		v.URL = App.MustReverse(templateDownloadHandlerName, v.Name, v.Version)
	}
}

func main() {
	App.MustListenAndServe()
}
