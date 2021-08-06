package main

import (
	"io"
	"net/http"
	"text/template"

	"github.com/fasibio/micropuzzle/logger"
	"github.com/go-chi/chi"
)

func main() {
	logger.Initialize("info")
	r := chi.NewRouter()
	ChiFileServer(r, "/", http.Dir("./public"))
	logger.Get().Infow("Start Server on Port :3000")
	logger.Get().Fatal(http.ListenAndServe(":3000", r))
}

func ChiFileServer(r chi.Router, path string, root http.FileSystem) {
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}
		f, err := root.Open(path)

		if err == nil {
			err := handleTemplate(f, w)
			if err != nil {
				logger.Get().Errorw("Error handle template", "error", err)
			}
		} else {
			logger.Get().Info("Will return Fallback ", path)
			f, err = root.Open("/index.html")
			handleTemplate(f, w)
			if err != nil {
				logger.Get().Error("Error by return fallback ", err)
			}
		}
	})
}

func handleTemplate(f http.File, dst io.Writer) error {
	text, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	handler := NewTemplateHandler()
	tmpl, err := template.New("httptemplate").Parse(string(text))
	if err != nil {
		return err
	}
	return tmpl.Execute(dst, handler)
}
