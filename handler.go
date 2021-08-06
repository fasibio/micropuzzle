package main

import (
	"fmt"

	"github.com/fasibio/micropuzzle/logger"
)

func NewTemplateHandler() TemplateHandler {
	return TemplateHandler{
		Reader: Reader{
			Test: "Test123",
		},
	}
}

type TemplateHandler struct {
	Reader Reader
}

type Reader struct {
	Test string
}

func (r Reader) Load(url string) string {
	logger.Get().Infow("load", "dest", url)
	return fmt.Sprintf("<a href='%s'>NANA</a>", url)
}
