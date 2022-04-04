package templatehandling

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fasibio/micropuzzle/configloader"
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
	frontends configloader.Configuration
}

func NewTemplateHandler(r *http.Request, socketUrl string, id uuid.UUID, server FragmentHandling, frontends configloader.Configuration, page configloader.Page) (*TemplateHandler, error) {
	return &TemplateHandler{
		socketUrl: socketUrl,
		Reader:    NewReader(server, r, id, frontends, page),
		frontends: frontends,
	}, nil
}

func (t *TemplateHandler) ScriptLoader() string {
	return "<script type=\"module\" src=\"/micro-lib/micropuzzle-components.esm.js\"></script>"
}

func (t *TemplateHandler) Loader() string {
	pagesbytes, _ := json.Marshal(t.frontends.GetPagesList())
	return fmt.Sprintf("<micro-puzzle-loader pagesStr='%s' streamingUrl=\"%s\" streamRegisterName=\"%s\" fallbacks=\"%d\"></micro-puzzle-loader>", string(pagesbytes), t.socketUrl, t.Reader.GetRequestId(), t.Reader.GetFallbacks())
}
