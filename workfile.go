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
		"assets",
		"data",
		"tmpl",
	}
)

func Clean(s *wb.Session) {
	s.LRun("gondola clean")
}

func Build(s *wb.Session) {
	s.LRun("go build -ldflags '-r lib'")
}

func BeforeDeploy(s *wb.Session) {
	Clean(s)
	Build(s)
}

func Deploy(s *wb.Session) {
	s.Install(deps)
	app, _ := s.StringVar("app")
	s.PushDir(app)
	s.Upload(app, "")
	g := governator.New(s)
	__ = g.Stop(app)
	s.Sync(dirs...)
	s.UploadTemplate("app.conf", "")
	g.AddService(app, nil, app)
	g.Start(app)
	n, _ := nginx.New(s)
	port, _ := s.StringVar("port")
	n.AddServer(app, "127.0.0.1:"+port)
	n.AddSite(app, wb.Template("conf/nginx/site.conf"))
	n.StartSite(app)
}

func init() {
	wb.DefaultJob = "Deploy"
	wb.SetDefaultHosts("governator.io")
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	wb.Var("app", filepath.Base(dir), "Application binary")
	wb.AddVars("app.conf")
}
