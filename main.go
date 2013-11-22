package main

import (
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

func main() {
	config.MustParse(&Config)
	m := mux.New()
	if !admin.Perform(m) {
		m.HandleFunc("^/$", mux.TemplateHandler("main.html", util.M{"Section": "home"}))
		m.SetTrustXHeaders(true)
		m.HandleAssets("/static/", STATIC_FILES_PATH)
		m.MustListenAndServe()
	}
}
