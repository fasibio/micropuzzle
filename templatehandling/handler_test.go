package templatehandling

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTemplateHandler_ScriptLoader(t *testing.T) {
	handler := TemplateHandler{}
	assert.Equal(t, handler.ScriptLoader(), "<script type=\"module\" src=\"/micro-lib/micropuzzle-components.esm.js\"></script>")
}

type ReaderMock struct {
	mock.Mock
}

func (r *ReaderMock) Load(url string, content string) string {
	args := r.Called(url, content)
	return args.String(0)
}

func (r *ReaderMock) GetRequestId() uuid.UUID {
	args := r.Called()
	return args.Get(0).(uuid.UUID)
}

func (r *ReaderMock) GetFallbacks() int64 {
	args := r.Called()
	return args.Get(0).(int64)
}

func TestTemplateHandler_Loader(t *testing.T) {

	readerMock := new(ReaderMock)
	id := uuid.Must(uuid.NewV4())

	readerMock.On("GetRequestId").Return(id)
	readerMock.On("GetFallbacks").Return(int64(0))
	handler := TemplateHandler{
		socketUrl: "socket_url",
		Reader:    readerMock,
	}
	assert.Equal(t, handler.Loader(), "<micro-puzzle-loader streamingUrl=\"socket_url\" streamRegisterName=\""+id.String()+"\" fallbacks=\"0\"></micro-puzzle-loader>")
}
