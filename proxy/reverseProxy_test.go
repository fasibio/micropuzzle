package proxy

import (
	"net/http/httptest"
	"testing"

	"github.com/fasibio/micropuzzle/configloader"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/suite"
)

type ReverseProxyTestSuite struct {
	suite.Suite
}

func TestReverseProxy(t *testing.T) {
	suite.Run(t, new(ReverseProxyTestSuite))
}

func (s *ReverseProxyTestSuite) TestRegisterReverseProxy() {
	s.Run("Happy Path", func() {
		r := chi.NewRouter()
		frontends := make(configloader.Frontends)
		child := make(map[string]configloader.Frontend)
		child["child"] = configloader.Frontend{
			Url: "http://test.url",
		}
		frontends["test"] = child
		RegisterReverseProxy(r, frontends)
		ts := httptest.NewServer(r)
		defer ts.Close()
		// TODO finish test
	})
}
