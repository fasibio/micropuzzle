package templatehandling

import (
	"fmt"
	"net/http"

	"github.com/fasibio/micropuzzle/proxy"
	"github.com/gofrs/uuid"
)

type ReaderHandling interface {
	GetRequestId() uuid.UUID
	GetFallbacks() int64
}

type FragmentHandling interface {
	LoadFragment(frontend, fragmentName, userId, remoteAddr string, header http.Header) (string, proxy.CacheInformation, bool)
}

type TemplateHandler struct {
	socketUrl string
	Reader    ReaderHandling
}

func NewTemplateHandler(r *http.Request, socketUrl string, id uuid.UUID, server FragmentHandling) (*TemplateHandler, error) {
	return &TemplateHandler{
		socketUrl: socketUrl,
		Reader:    NewReader(server, r, id),
	}, nil
}

func (t *TemplateHandler) ScriptLoader() string {
	return "<script type=\"module\" src=\"/micro-lib/micropuzzle-components.esm.js\"></script>"
}

func (t *TemplateHandler) Loader() string {
	return fmt.Sprintf("<micro-puzzle-loader streamingUrl=\"%s\" streamRegisterName=\"%s\" fallbacks=\"%d\"></micro-puzzle-loader>", t.socketUrl, t.Reader.GetRequestId(), t.Reader.GetFallbacks())
}
