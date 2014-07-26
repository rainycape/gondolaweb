package main

import (
	"path"

	"gnd.la/app"
	"gnd.la/apps/docs"
	"gnd.la/net/urlutil"
)

const (
	//gondolaURL = "http://www.gondolaweb.com"
	gondolaURL = "ssh://abra.rm-fr.net/home/fiam/git/gondola.git"
)

func gndlaHandler(ctx *app.Context) {
	if ctx.FormValue("go-get") == "1" {
		ctx.MustExecute("goget.html", nil)
		return
	}
	// Check if the request path is a pkg name
	var p string
	pkg := path.Join("gnd.la", ctx.R.URL.Path)
	if _, err := docs.DefaultContext.Import(pkg, "", 0); err == nil {
		p = ctx.MustReverse(docs.PackageHandlerName, pkg)
	}
	redir, err := urlutil.Join(gondolaURL, p)
	if err != nil {
		panic(err)
	}
	ctx.Redirect(redir, false)
}
