package main

import (
	"gnd.la/app"
	"gnd.la/config"
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
	App.HandleAssets("/assets/", pathutil.Relative("assets"))
	// You might probably want the following if you're
	// deploying your app behind an upstream proxy.
	//
	// App.SetTrustXHeaders(true)

	// Site handlers
	App.Handle("^/$", MainHandler, app.NamedHandler("main"))
}
