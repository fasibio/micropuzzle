package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fasibio/micropuzzle/logger"
	"github.com/fasibio/micropuzzle/proxy"
	"github.com/gofrs/uuid"
)

func NewTemplateHandler(r *http.Request, timeout time.Duration, cache ChacheHandler) (*TemplateHandler, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	return &TemplateHandler{
		Loader: fmt.Sprintf("<micro-puzzle-loader streamRegisterName=\"%s\"></micro-puzzle-loader>", id),
		Reader: Reader{
			requestId:   id,
			cache:       cache,
			timeout:     timeout,
			mainRequest: r,
			proxy:       proxy.Proxy{},
		},
	}, nil
}

type TemplateHandler struct {
	Reader Reader
	Loader string
}

type Reader struct {
	cache       ChacheHandler
	timeout     time.Duration
	proxy       proxy.Proxy
	mainRequest *http.Request
	requestId   uuid.UUID
}

func (r *Reader) Load(url, content string) string {
	log.Println("nanan", content)
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
			return "<h1>Fallback</h1>"
		}
	}
}

func (r *Reader) loadAsync(url string, content string, result *chan string, timeout *chan bool) {
	logger.Get().Infow("load", "dest", url)
	res, err := r.proxy.Get(url, r.mainRequest)

	if err != nil {
		logger.Get().Warnw("error by load url", "url", url, "error", err)
		return
	}

	log.Println("result", content, len(*timeout))
	contentPage := fmt.Sprintf("<micro-puzzle-element name=\"%s\"><template>%s</template></micro-puzzle-element>", content, string(res))
	if len(*timeout) == 1 {
		err := r.cache.add(fmt.Sprintf("%v_%s", r.requestId, content), []byte(contentPage))
		if err != nil {
			logger.Get().Warnw("error by saving to cache", "error", err)
		}
	} else {
		*result <- contentPage
	}
}
