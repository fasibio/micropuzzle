package templatehandling

import (
	"fmt"
	"net/http"

	"github.com/fasibio/micropuzzle/fragments"
	"github.com/gofrs/uuid"
)

func NewTemplateHandler(r *http.Request, socketUrl string, id uuid.UUID, server *fragments.FragmentHandler) (*TemplateHandler, error) {

	return &TemplateHandler{
		socketUrl: socketUrl,
		Reader: Reader{
			hasFallbacks: 0,
			requestId:    id,
			mainRequest:  r,
			server:       server,
		},
	}, nil
}

type TemplateHandler struct {
	socketUrl string
	Reader    Reader
}

func (t *TemplateHandler) ScriptLoader() string {
	return "<script type=\"module\" src=\"/micro-lib/micropuzzle-components.esm.js\"></script>"
}

func (t *TemplateHandler) Loader() string {
	return fmt.Sprintf("<micro-puzzle-loader streamingUrl=\"%s\" streamRegisterName=\"%s\" fallbacks=\"%d\"></micro-puzzle-loader>", t.socketUrl, t.Reader.requestId, t.Reader.hasFallbacks)
}
