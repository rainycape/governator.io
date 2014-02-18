// +build NONE

package main

import (
	"os"
	"path/filepath"
	"workbot.io/contrib/governator"
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
	app, _ := s.StringVar("App")
	s.PushDir(app)
	s.UploadFile(app, "")
	g := governator.New(s)
	__ = g.Stop(app)
	s.Sync(dirs...)
	s.UploadTemplateFile("conf/governator/app", "/etc/governator/"+app)
	s.UploadTemplateFile("conf/nginx/nginx.conf", "/etc/nginx/conf.d/"+app)
	s.UploadTemplateFile("conf/nginx/site.conf", "/etc/nginx/"+app)
	s.UploadTemplateFile("app.conf", "")
	g.Start(app)
}

func init() {
	wb.DefaultJob = "Deploy"
	wb.SetDefaultHosts("governator.io")
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	wb.Var("App", filepath.Base(dir), "Application binary")
	wb.AddContextFile("app.conf")
}
