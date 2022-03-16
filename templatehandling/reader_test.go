package templatehandling

import (
	"net/http"
	"testing"

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
			r := &Reader{}
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
	frontend := "frontend"
	fragmentName := "fragmentName"
	remoteAddr := "remoteAddr"
	httpHeaderMock := http.Header{
		"User-Agent": []string{"userAgent"},
	}
	id := uuid.Must(uuid.NewV4())

	s.Run("Happy Path", func() {
		mockObj := new(MockedFragmentHandler)
		mockObj.On("LoadFragment", frontend, fragmentName, id.String(), remoteAddr, httpHeaderMock).Return("result", proxy.CacheInformation{}, false)
		reader := Reader{
			server:    mockObj,
			requestId: id,
			mainRequest: &http.Request{
				RemoteAddr: remoteAddr,
				Header:     httpHeaderMock,
			},
		}
		result := reader.Load(frontend, fragmentName)
		assert.Equal(s.T(), result, "<micro-puzzle-element name=\"fragmentName\"><template>result</template></micro-puzzle-element>")
		mockObj.AssertExpectations(s.T())
	})
	s.Run("Loading need To long so fallback will be returned", func() {
		mockObj := new(MockedFragmentHandler)
		mockObj.On("LoadFragment", frontend, fragmentName, id.String(), remoteAddr, httpHeaderMock).Return("fallbackhtml", proxy.CacheInformation{}, true)
		reader := Reader{
			server:       mockObj,
			requestId:    id,
			hasFallbacks: 0,
			mainRequest: &http.Request{
				RemoteAddr: remoteAddr,
				Header:     httpHeaderMock,
			},
		}
		result := reader.Load(frontend, fragmentName)
		assert.Equal(s.T(), result, "<micro-puzzle-element name=\"fragmentName\"><template>fallbackhtml</template></micro-puzzle-element>")
		mockObj.AssertExpectations(s.T())
		assert.Equal(s.T(), reader.hasFallbacks, int64(1))
	})

}
