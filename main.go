package main

import (
	"bytes"
	"fmt"
	"gnd.la/app"
	_ "gnd.la/bootstrap"
	"gnd.la/config"
	"gnd.la/html"
	"gnd.la/template"
	"gnd.la/util"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	Config config.Config
	App    *app.App
	data   = util.RelativePath("data")
)

func fileHandler(name string) app.Handler {
	p := filepath.Join(data, name)
	return func(ctx *app.Context) {
		f, err := os.Open(p)
		if err != nil {
			ctx.NotFound("")
		}
		defer f.Close()
		var modtime time.Time
		if fi, err := f.Stat(); err == nil {
			modtime = fi.ModTime()
		}
		http.ServeContent(ctx, ctx.R, name, modtime, f)
	}
}

func init() {
	template.AddFuncs(template.FuncMap{
		"#config": func(name string, def string) template.HTML {
			var buf bytes.Buffer
			buf.WriteString("<h3 class=\"config\">")
			buf.WriteString(html.Escape(name))
			if def != "" {
				buf.WriteString(" <span class=\"label label-success\">optional</span>")
				fmt.Fprintf(&buf, " <span class=\"default\">default: %s</span>", html.Escape(def))
			} else {
				buf.WriteString(" <span class=\"label label-danger\">required</span>")
			}
			buf.WriteString("</h3>")
			return template.HTML(buf.String())
		}})
	config.MustParse(&Config)
	App = app.New()
	App.SetTrustXHeaders(true)
	App.HandleAssets("/assets/", util.RelativePath("assets"))
	App.Handle("^/$", app.TemplateHandler("main.html", nil))
	App.Handle("^/install\\.sh$", fileHandler("install.sh"))
	App.Handle("^/get/releases/linux/x86_64/latest/governator$", fileHandler("governator"))
	App.Handle("^/contrib/upstart/governator.conf$", fileHandler("governator.conf"))
}

func main() {
	App.MustListenAndServe()
}
