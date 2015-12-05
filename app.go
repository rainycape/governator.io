package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gnd.la/app"
	"gnd.la/config"
	_ "gnd.la/frontend/bootstrap"
	"gnd.la/html"
	"gnd.la/template"
	"gnd.la/util/pathutil"
)

var (
	App  *app.App
	data = pathutil.Relative("data")
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
	config.MustParse()
	App = app.New()
	App.SetTrustXHeaders(true)

	// Redirect all other possible hosts to governator.io
	redir := app.RedirectHandler("http://governator.io${0}", true)
	App.Handle("(.*)", redir, app.HostHandler("governator-io.appspot.com"))
	App.Handle("(.*)", redir, app.HostHandler("www.governator.io"))

	App.HandleAssets("/assets/", pathutil.Relative("assets"))
	App.Handle("^/$", app.TemplateHandler("main.html", nil))
	App.Handle("^/install\\.sh$", fileHandler("contrib/install.sh"))
	App.Handle("^/get/releases/linux/x86_64/latest/governator$", fileHandler("governator"))

	App.HandleAssets("/contrib/", pathutil.Relative("data/contrib"))
}
