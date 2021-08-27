package main

import (
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
)

func NewTemplateHandler(r *http.Request, socketUrl string, id uuid.UUID, server *SocketHandler) (*TemplateHandler, error) {

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
	server      *SocketHandler
	mainRequest *http.Request
	requestId   uuid.UUID
}

func (r *Reader) Load(url, content string) string {
	return r.server.Load(url, content, r.requestId, r.mainRequest.Header, r.mainRequest.RemoteAddr)
}

type NewContentPayload struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}
