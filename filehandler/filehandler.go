package filehandler

import (
	"io"
	"net/http"
	"text/template"

	"github.com/fasibio/micropuzzle/logger"
	"github.com/fasibio/micropuzzle/proxy"
	"github.com/fasibio/micropuzzle/templatehandling"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

type FragmentHandling interface {
	LoadFragment(frontend, fragmentName, userId, remoteAddr string, header http.Header) (string, proxy.CacheInformation, bool)
}

type templateCreator func(r *http.Request, socketUrl string, id uuid.UUID, server templatehandling.FragmentHandling) (*templatehandling.TemplateHandler, error)

type fileHandler struct {
	Server          FragmentHandling
	SocketUrl       string
	templateHandler templateCreator
}

func NewFileHandler(server FragmentHandling, socketUrl string) *fileHandler {
	return &fileHandler{
		Server:          server,
		SocketUrl:       socketUrl,
		templateHandler: templatehandling.NewTemplateHandler,
	}
}

func (filehandler *fileHandler) RegisterFileHandler(r chi.Router, path string, root http.FileSystem) {
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

func (filehandler *fileHandler) handleTemplate(f http.File, dst io.Writer, r *http.Request) error {
	text, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	handler, err := filehandler.templateHandler(r, filehandler.SocketUrl, id, filehandler.Server)
	if err != nil {
		return err
	}
	tmpl, err := template.New("httptemplate").Parse(string(text))
	if err != nil {
		return err
	}

	return tmpl.Execute(dst, handler)
}
