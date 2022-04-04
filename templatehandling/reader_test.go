package templatehandling

import (
	"net/http"
	"testing"

	"github.com/fasibio/micropuzzle/configloader"
	"github.com/fasibio/micropuzzle/proxy"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestReader_getMicroPuzzleElement(t *testing.T) {
	type args struct {
		name    string
		content string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test result by given params",
			args: args{
				name:    "test_name",
				content: "test:content",
			},
			want: "<micro-puzzle-element name=\"test_name\"><template>test:content</template></micro-puzzle-element>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &reader{}
			if got := r.getMicroPuzzleElement(tt.args.name, tt.args.content); got != tt.want {
				t.Errorf("Reader.getMicroPuzzleElement() = %v, want %v", got, tt.want)
			}
		})
	}
}

type MockedFragmentHandler struct {
	mock.Mock
}

func (m MockedFragmentHandler) LoadFragment(frontend, fragmentName, userId, remoteAddr string, header http.Header) (string, proxy.CacheInformation, bool) {
	args := m.Called(frontend, fragmentName, userId, remoteAddr, header)
	return args.String(0), args[1].(proxy.CacheInformation), args.Bool(2)
}

type ReaderTestSuite struct {
	suite.Suite
}

func TestReaderTestSuit(t *testing.T) {
	suite.Run(t, new(ReaderTestSuite))
}

func (s ReaderTestSuite) TestReader_Load() {
	remoteAddr := "remoteAddr"
	httpHeaderMock := http.Header{
		"User-Agent": []string{"userAgent"},
	}
	area := "header"
	id := uuid.Must(uuid.NewV4())
	s.Run("Happy Path", func() {
		mockObj := new(MockedFragmentHandler)
		mockObj.On("LoadFragment", "global.header", area, id.String(), remoteAddr, httpHeaderMock).Return("result", proxy.CacheInformation{}, false)
		reader := NewReader(mockObj, &http.Request{
			RemoteAddr: remoteAddr,
			Header:     httpHeaderMock,
		}, id, configloader.Configuration{
			Definitions: map[string]map[string]configloader.Definition{"global": {"header": configloader.Definition{Url: "mockurl"}}},
		}, configloader.Page{Url: "/", Fragments: map[string]string{"header": "global.header"}})
		result := reader.Load(area)
		assert.Equal(s.T(), result, "<micro-puzzle-element name=\"header\"><template>result</template></micro-puzzle-element>")
		mockObj.AssertExpectations(s.T())
	})

	s.Run("Loading need To long so fallback will be returned", func() {
		mockObj := new(MockedFragmentHandler)
		mockObj.On("LoadFragment", "global.header", area, id.String(), remoteAddr, httpHeaderMock).Return("fallback", proxy.CacheInformation{}, true)
		reader := NewReader(mockObj, &http.Request{
			RemoteAddr: remoteAddr,
			Header:     httpHeaderMock,
		}, id, configloader.Configuration{
			Definitions: map[string]map[string]configloader.Definition{"global": {"header": configloader.Definition{Url: "mockurl"}}},
		}, configloader.Page{Url: "/", Fragments: map[string]string{"header": "global.header"}})
		result := reader.Load(area)
		assert.Equal(s.T(), result, "<micro-puzzle-element name=\"header\"><template>fallback</template></micro-puzzle-element>")
		mockObj.AssertExpectations(s.T())
		s.Assert().Equal(reader.fallbacks, int64(1))
	})
}
