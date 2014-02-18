package main

import (
	"gnd.la/app"
	_ "gnd.la/bootstrap"
	"gnd.la/config"
	"gnd.la/util"
)

var (
	Config config.Config
	App    *app.App
)

func init() {
	config.MustParse(&Config)
	App = app.New()
	App.SetTrustXHeaders(true)
	App.HandleAssets("/assets/", util.RelativePath("assets"))
	App.Handle("^/$", app.TemplateHandler("main.html", nil))
}

func main() {
	App.MustListenAndServe()
}
