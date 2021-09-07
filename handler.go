package main

import (
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
)

func NewTemplateHandler(r *http.Request, socketUrl string, id uuid.UUID, server *WebSocketHandler) (*TemplateHandler, error) {

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

type Reader struct {
	server       *WebSocketHandler
	mainRequest  *http.Request
	requestId    uuid.UUID
	hasFallbacks int64
}

func (r *Reader) Load(url, content string) string {
	result, isFallback := r.server.LoadFragment(url, content, r.requestId.String(), r.mainRequest.RemoteAddr, r.mainRequest.Header)
	if isFallback {
		r.hasFallbacks = r.hasFallbacks + 1
	}
	return r.getMicroPuzzleElement(content, result)
}

func (r *Reader) getMicroPuzzleElement(name, content string) string {
	return fmt.Sprintf("<micro-puzzle-element name=\"%s\"><template>%s</template></micro-puzzle-element>", name, content)
}
