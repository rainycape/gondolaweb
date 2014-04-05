package main

import (
	"gnd.la/app"
	"gnd.la/apps/docs"
	"gnd.la/apps/docs/doc"
	"gnd.la/util/urlutil"
	"path"
)

const (
	gondolaURL = "http://www.gondolaweb.com"
)

func gndlaHandler(ctx *app.Context) {
	if ctx.FormValue("go-get") == "1" {
		ctx.MustExecute("goget.html", nil)
		return
	}
	// Check if the request path is a pkg name
	var p string
	pkg := path.Join("gnd.la", ctx.R.URL.Path)
	if _, err := doc.Context.Import(pkg, "", 0); err == nil {
		p = ctx.MustReverse(docs.PackageHandlerName, pkg)
	}
	redir, err := urlutil.Join(gondolaURL, p)
	if err != nil {
		panic(err)
	}
	ctx.Redirect(redir, false)
}
