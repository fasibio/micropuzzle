package filehandler

import (
	"io"
	"net/http"
	"text/template"

	"github.com/fasibio/micropuzzle/configloader"
	"github.com/fasibio/micropuzzle/logger"
	"github.com/fasibio/micropuzzle/mimetypefinder"
	"github.com/fasibio/micropuzzle/proxy"
	"github.com/fasibio/micropuzzle/templatehandling"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

type FragmentHandling interface {
	LoadFragment(frontend, fragmentName, userId, remoteAddr string, header http.Header) (string, proxy.CacheInformation, bool)
}

type templateCreator func(r *http.Request, socketUrl string, id uuid.UUID, server templatehandling.FragmentHandling, frontends configloader.Configuration, page configloader.Page) (*templatehandling.TemplateHandler, error)

type fileHandler struct {
	Server          FragmentHandling
	SocketUrl       string
	templateHandler templateCreator
	config          configloader.Configuration
}

func NewFileHandler(server FragmentHandling, socketUrl string, config configloader.Configuration) *fileHandler {
	return &fileHandler{
		Server:          server,
		SocketUrl:       socketUrl,
		templateHandler: templatehandling.NewTemplateHandler,
		config:          config,
	}
}

func (filehandler *fileHandler) RegisterFileHandler(r chi.Router, root http.FileSystem) {
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}
		f, err := root.Open(path)
		mimetype := mimetypefinder.MimeTypeForFile(path)
		w.Header().Set("Content-Type", mimetype)
		if err == nil {
			if mimetype == "application/javascript" {
				io.Copy(w, f)
				return
			}
			err := filehandler.handleTemplate(f, w, r, *filehandler.config.Pages.GetPageByUrl("/"))
			if err != nil {
				logger.Get().Warnw("Error handle template", "error", err)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

func (filehandler *fileHandler) HandlePage(r chi.Router, p configloader.Page, root http.FileSystem) {
	logger.Get().Infow("Register page endpoint", "path", p.Url)
	r.Get(p.Url, func(w http.ResponseWriter, r *http.Request) {
		path := p.Template
		f, err := root.Open(path)
		if err != nil {
			logger.Get().Warnw("Error open file", "error", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		err = filehandler.handleTemplate(f, w, r, p)
		if err != nil {
			logger.Get().Warnw("Error handle template", "error", err)
		}
	})
}

func (filehandler *fileHandler) handleTemplate(f io.Reader, dst io.Writer, r *http.Request, p configloader.Page) error {
	text, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	handler, err := filehandler.templateHandler(r, filehandler.SocketUrl, id, filehandler.Server, filehandler.config, p)
	if err != nil {
		return err
	}
	tmpl, err := template.New("httptemplate").Parse(string(text))
	if err != nil {
		return err
	}

	return tmpl.Execute(dst, handler)
}
