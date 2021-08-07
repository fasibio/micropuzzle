package main

import (
	"net/http"

	"github.com/fasibio/micropuzzle/logger"
	"github.com/fasibio/micropuzzle/proxy"
)

func NewTemplateHandler(r *http.Request) TemplateHandler {
	return TemplateHandler{
		Reader: Reader{
			Test:        "Test123",
			mainRequest: r,
			proxy:       proxy.Proxy{},
		},
	}
}

type TemplateHandler struct {
	Reader Reader
}

type Reader struct {
	Test        string
	proxy       proxy.Proxy
	mainRequest *http.Request
}

func (r Reader) Load(url string) string {
	logger.Get().Infow("load", "dest", url)
	result, err := r.proxy.Get(url, r.mainRequest)
	if err != nil {
		logger.Get().Warnw("error by load url", "url", url, "error", err)
	}
	return string(result)
}
