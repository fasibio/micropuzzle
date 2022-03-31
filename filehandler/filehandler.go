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

type templateCreator func(r *http.Request, socketUrl string, id uuid.UUID, server templatehandling.FragmentHandling, frontends configloader.Frontends) (*templatehandling.TemplateHandler, error)

type fileHandler struct {
	Server          FragmentHandling
	SocketUrl       string
	templateHandler templateCreator
	config          configloader.Frontends
}

func NewFileHandler(server FragmentHandling, socketUrl string, config configloader.Frontends) *fileHandler {
	return &fileHandler{
		Server:          server,
		SocketUrl:       socketUrl,
		templateHandler: templatehandling.NewTemplateHandler,
		config:          config,
	}
}

func (filehandler *fileHandler) RegisterFileHandler(r chi.Router, path string, root http.FileSystem) {
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
			err := filehandler.handleTemplate(f, w, r)
			if err != nil {
				logger.Get().Warnw("Error handle template", "error", err)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

func (filehandler *fileHandler) HandlePage(r chi.Router, p configloader.Page, root http.FileSystem) {
	r.Get(p.Url, func(w http.ResponseWriter, r *http.Request) {
		f, err := root.Open("/index.html") // TODO add this to page config
		if err != nil {
			logger.Get().Warnw("Error open file", "error", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		err = filehandler.handleTemplateV2(f, w, r)
		if err != nil {
			logger.Get().Warnw("Error handle template", "error", err)
		}
	})
}

func (filehandler *fileHandler) handleTemplateV2(f io.Reader, dst io.Writer, r *http.Request, p configloader.Page) error {
	text, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	handler, err := filehandler.templateHandler(r, filehandler.SocketUrl, id, filehandler.Server, filehandler.config)
	if err != nil {
		return err
	}
	tmpl, err := template.New("httptemplate").Parse(string(text))
	if err != nil {
		return err
	}

	return tmpl.Execute(dst, handler)
}

func (filehandler *fileHandler) handleTemplate(f io.Reader, dst io.Writer, r *http.Request) error {
	text, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	handler, err := filehandler.templateHandler(r, filehandler.SocketUrl, id, filehandler.Server, filehandler.config)
	if err != nil {
		return err
	}
	tmpl, err := template.New("httptemplate").Parse(string(text))
	if err != nil {
		return err
	}

	return tmpl.Execute(dst, handler)
}
