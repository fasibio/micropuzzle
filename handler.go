package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fasibio/micropuzzle/logger"
	"github.com/fasibio/micropuzzle/proxy"
	"github.com/gofrs/uuid"
	socketio "github.com/googollee/go-socket.io"
)

func NewTemplateHandler(r *http.Request, timeout time.Duration, cache ChacheHandler, socketUrl string, id uuid.UUID, server *socketio.Server) (*TemplateHandler, error) {

	return &TemplateHandler{
		Loader: fmt.Sprintf("<micro-puzzle-loader streamingUrl=\"%s\" streamRegisterName=\"%s\"></micro-puzzle-loader>", socketUrl, id),
		Reader: Reader{
			requestId:   id,
			cache:       cache,
			timeout:     timeout,
			mainRequest: r,
			proxy:       proxy.Proxy{},
			server:      server,
		},
	}, nil
}

type TemplateHandler struct {
	Reader Reader
	Loader string
}

type Reader struct {
	server      *socketio.Server
	cache       ChacheHandler
	timeout     time.Duration
	proxy       proxy.Proxy
	mainRequest *http.Request
	requestId   uuid.UUID
}

func (r *Reader) Load(url, content string) string {
	resultChan := make(chan string, 1)
	timeout := make(chan bool, 1)
	timeoutBubble := make(chan bool, 1)
	go r.loadAsync(url, content, &resultChan, &timeoutBubble)

	go func() {
		time.Sleep(r.timeout)
		timeout <- true
	}()
	select {
	case d := <-resultChan:
		{
			return d
		}
	case <-timeout:
		{
			timeoutBubble <- true
			return fmt.Sprintf("<micro-puzzle-element name=\"%s\"><template>%s</template></micro-puzzle-element>", content, "<h1>Fallback</h1>")
		}
	}
}

type NewContentPayload struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

func (r *Reader) loadAsync(url string, content string, result *chan string, timeout *chan bool) {
	res, err := r.proxy.Get(url, r.mainRequest)

	if err != nil {
		logger.Get().Warnw("error by load url", "url", url, "error", err)
		return
	}

	contentPage := fmt.Sprintf("<micro-puzzle-element name=\"%s\"><template>%s</template></micro-puzzle-element>", content, string(res))
	if len(*timeout) == 1 {
		if r.server.RoomLen("", r.requestId.String()) > 0 {
			r.server.BroadcastToRoom("/", r.requestId.String(), "NEW_CONTENT", NewContentPayload{Key: content, Value: string(res)})
		} else {
			err := r.cache.Add(r.requestId.String(), content, []byte(contentPage))
			if err != nil {
				logger.Get().Warnw("error by saving to cache", "error", err)
			}
		}
	} else {
		*result <- contentPage
	}
}
