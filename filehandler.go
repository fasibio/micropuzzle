package main

import (
	"io"
	"net/http"
	"text/template"

	"github.com/fasibio/micropuzzle/fragments"
	"github.com/fasibio/micropuzzle/logger"
	"github.com/fasibio/micropuzzle/templatehandling"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

type FileHandler struct {
	server    *fragments.FragmentHandler
	socketUrl string
}

func (filehandler *FileHandler) ChiFileServer(r chi.Router, path string, root http.FileSystem) {

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}
		f, err := root.Open(path)
		mimetype := mimeTypeForFile(path)
		w.Header().Set("Content-Type", mimetype)
		if err == nil {
			if mimetype == "application/javascript" {
				io.Copy(w, f)
				return
			}
			err := filehandler.handleTemplate(f, w, r)
			if err != nil {
				logger.Get().Warnw("Error handle template", "error", err)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

func (filehandler *FileHandler) handleTemplate(f http.File, dst io.Writer, r *http.Request) error {

	text, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	handler, err := templatehandling.NewTemplateHandler(r, filehandler.socketUrl, id, filehandler.server)
	if err != nil {
		return err
	}
	tmpl, err := template.New("httptemplate").Parse(string(text))
	if err != nil {
		return err
	}

	return tmpl.Execute(dst, handler)
}
