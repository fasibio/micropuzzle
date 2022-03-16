package templatehandling

import (
	"fmt"
	"net/http"

	"github.com/fasibio/micropuzzle/fragments"
	"github.com/gofrs/uuid"
)

type Reader struct {
	server       fragments.FragmentHandling
	mainRequest  *http.Request
	requestId    uuid.UUID
	hasFallbacks int64
}

func (r *Reader) Load(url, content string) string {
	result, _, isFallback := r.server.LoadFragment(url, content, r.requestId.String(), r.mainRequest.RemoteAddr, r.mainRequest.Header)
	if isFallback {
		r.hasFallbacks = r.hasFallbacks + 1
	}

	return r.getMicroPuzzleElement(content, result)
}

func (r *Reader) getMicroPuzzleElement(name, content string) string {
	return fmt.Sprintf("<micro-puzzle-element name=\"%s\"><template>%s</template></micro-puzzle-element>", name, content)
}
