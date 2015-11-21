package main

import (
	"gnd.la/app"
	"gnd.la/config"
	"gnd.la/template/assets"
	"gnd.la/util/pathutil"
)

var (
	App *app.App
)

var Config struct {
	// Include your settings here
}

func init() {
	// Initialize the configuration and the App in init, so
	// it's configured correctly when running tests.
	config.Register(&Config)
	config.MustParse()

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
	// App.SetTrustXHeaders(true)

	// Site handlers
	App.Handle("^/$", MainHandler, app.NamedHandler("main"))

	// Error handler, for 404
	App.SetErrorHandler(ErrorHandler)
}
