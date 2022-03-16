package templatehandling

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplateHandler_ScriptLoader(t *testing.T) {
	handler := TemplateHandler{}
	assert.Equal(t, handler.ScriptLoader(), "<script type=\"module\" src=\"/micro-lib/micropuzzle-components.esm.js\"></script>")
}

func TestTemplateHandler_Loader(t *testing.T) {
	s := [16]byte{65, 66, 67, 226, 130, 172}
	handler := TemplateHandler{
		socketUrl: "socket_url",
		Reader: Reader{
			requestId:    s,
			hasFallbacks: 0,
		},
	}
	assert.Equal(t, handler.Loader(), "<micro-puzzle-loader streamingUrl=\"socket_url\" streamRegisterName=\"414243e2-82ac-0000-0000-000000000000\" fallbacks=\"0\"></micro-puzzle-loader>")
}
