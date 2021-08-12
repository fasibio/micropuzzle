package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fasibio/micropuzzle/logger"
	"github.com/fasibio/micropuzzle/proxy"
)

func NewTemplateHandler(r *http.Request, timeout time.Duration) TemplateHandler {
	return TemplateHandler{
		Reader: Reader{
			timeout:     timeout,
			mainRequest: r,
			proxy:       proxy.Proxy{},
		},
	}
}

type TemplateHandler struct {
	Reader Reader
}

type Reader struct {
	timeout     time.Duration
	proxy       proxy.Proxy
	mainRequest *http.Request
}

func (r Reader) Load(url, content string) string {
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

func (r Reader) loadAsync(url string, content string, result *chan string, timeout *chan bool) {
	logger.Get().Infow("load", "dest", url)
	res, err := r.proxy.Get(url, r.mainRequest)

	if err != nil {
		logger.Get().Warnw("error by load url", "url", url, "error", err)
		return
	}

	log.Println("result", content, len(*timeout))
	if len(*timeout) == 1 {
		// @TODO SAVE to cache for streaming service
		log.Println("have save to cache")
	} else {
		*result <- fmt.Sprintf("<micro-puzzle-element name=\"%s\"><template>%s</template></micro-puzzle-element>", content, string(res))

	}
}
