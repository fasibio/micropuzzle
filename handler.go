package main

import (
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
)

func NewTemplateHandler(r *http.Request, socketUrl string, id uuid.UUID, server *WebSocketHandler) (*TemplateHandler, error) {

	return &TemplateHandler{
		Loader: fmt.Sprintf("<micro-puzzle-loader streamingUrl=\"%s\" streamRegisterName=\"%s\"></micro-puzzle-loader>", socketUrl, id),
		Reader: Reader{
			requestId:   id,
			mainRequest: r,
			server:      server,
		},
	}, nil
}

type TemplateHandler struct {
	Reader Reader
	Loader string
}

type Reader struct {
	server      *WebSocketHandler
	mainRequest *http.Request
	requestId   uuid.UUID
}

func (r *Reader) Load(url, content string) string {
	result := r.server.Load(url, content, r.requestId.String(), r.mainRequest.Header, r.mainRequest.RemoteAddr)
	return r.getMicroPuzzleElement(content, result)
}

func (r *Reader) getMicroPuzzleElement(name, content string) string {
	return fmt.Sprintf("<micro-puzzle-element name=\"%s\"><template>%s</template></micro-puzzle-element>", name, content)
}
