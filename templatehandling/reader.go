package templatehandling

import (
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
)

type reader struct {
	server      FragmentHandling
	mainRequest *http.Request
	requestId   uuid.UUID
	fallbacks   int64
}

func NewReader(server FragmentHandling, r *http.Request, id uuid.UUID) *reader {
	return &reader{
		server:      server,
		mainRequest: r,
		requestId:   id,
		fallbacks:   0,
	}
}

func (r *reader) GetRequestId() uuid.UUID {
	return r.requestId
}

func (r *reader) GetFallbacks() int64 {
	return r.fallbacks
}

func (r *reader) Load(url, content string) string {
	result, _, isFallback := r.server.LoadFragment(url, content, r.requestId.String(), r.mainRequest.RemoteAddr, r.mainRequest.Header)
	if isFallback {
		r.fallbacks = r.fallbacks + 1
	}

	return r.getMicroPuzzleElement(content, result)
}

func (r *reader) getMicroPuzzleElement(name, content string) string {
	return fmt.Sprintf("<micro-puzzle-element name=\"%s\"><template>%s</template></micro-puzzle-element>", name, content)
}
