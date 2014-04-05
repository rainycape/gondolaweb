package main

import (
	"gnd.la/app"
	"gnd.la/config"
	"gnd.la/template/assets"
	"gnd.la/util/pathutil"
	"net/http"
)

var (
	App *app.App
)

var Config struct {
	config.Config
}

func init() {
	// Initialize the configuration and the App in init, so
	// it's configured correctly when running tests.

	config.MustParse(&Config)
	App = app.New()
	// Make the config available to templates as @Config
	App.AddTemplateVars(map[string]interface{}{
		"Config": &Config,
	})
	// Asset handling
	App.HandleAssets("/assets/", pathutil.Relative("assets"))

	// Gondola automatically servers favicon.ico and robots.txt
	// from the root, other assets that we want to be served
	// from the root have to be manually configured.
	handler := assets.Handler(assets.NewManager(App.AssetsManager().Loader(), ""))
	App.Handle("/(humans.txt|crossdomain.xml)$", func(ctx *app.Context) { handler(ctx, ctx.R) })

	// You might probably want the following if you're
	// deploying your app behind an upstream proxy.
	//
	// a.SetTrustXHeaders(true)

	// Site handlers
	App.HandleNamed("^/$", MainHandler, "main")

	// Error handler, for 404
	App.SetErrorHandler(ErrorHandler)
}

func MainHandler(ctx *app.Context) {
	ctx.MustExecute("main.html", nil)
}

func ErrorHandler(ctx *app.Context, msg string, code int) bool {
	// Only send the 404 page if the error is a 404
	if code == http.StatusNotFound {
		ctx.MustExecute("404.html", nil)
		return true
	}
	// Otherwise, let the Gondola error processing take care of it
	return false
}

func main() {
	// Start listening on main(), so the app does not start
	// listening when running tests.
	App.MustListenAndServe()
}
