// +build NONE

package main

import (
	"os"
	"path/filepath"

	"workbot.io/contrib/governator"
	"workbot.io/contrib/nginx"
	"workbot.io/wb"
)

var (
	deps = []string{}
	dirs = []string{
		"articles",
		"data",
		"tutorials",
		"templates",
		"assets",
		"tmpl",
	}
)

func Clean(s *wb.Session) {
	s.LRun("gondola clean")
}

func Build(s *wb.Session) {
	s.LRun("go build")
}

func BeforeDeploy(s *wb.Session) {
	Clean(s)
	Build(s)
}

func Deploy(s *wb.Session) {
	s.Install(deps)
	app, _ := s.StringVar("App")
	s.PushDir(app)
	s.Upload(app, "")
	s.Upload("app.conf", "")
	port, _ := s.StringVar("port")
	server := "127.0.0.1:" + port
	opts := map[string]string{"watchdog": "get http://" + server, "env": "GONDOLA_ALLOW_SHORT_SECRET=1"}
	g := governator.New(s)
	g.Install()
	g.AddService(app, opts, app)
	__ = g.Stop(app)
	s.Sync(dirs...)
	g.Start(app)
	n, _ := nginx.New(s)
	n.AddConfig(app, wb.Template("conf/nginx/nginx.conf"))
	n.AddSite(app, wb.Template("conf/nginx/site.conf"))
	n.StartSite(app)
}

func init() {
	wb.DefaultJob = "Deploy"
	wb.SetDefaultHosts("milos.rainycape.com")
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	wb.Var("App", filepath.Base(dir), "Application binary")
	wb.AddVars("app.conf")
}
