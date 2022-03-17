package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/fasibio/micropuzzle/configloader"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/suite"
)

type closeNotifyingRecorder struct {
	*httptest.ResponseRecorder
	closed chan bool
}

func newCloseNotifyingRecorder() *closeNotifyingRecorder {
	return &closeNotifyingRecorder{
		httptest.NewRecorder(),
		make(chan bool, 1),
	}
}

func (c *closeNotifyingRecorder) close() {
	c.closed <- true
}

func (c *closeNotifyingRecorder) CloseNotify() <-chan bool {
	return c.closed
}

type ReverseProxyTestSuite struct {
	suite.Suite
}

func TestReverseProxy(t *testing.T) {
	suite.Run(t, new(ReverseProxyTestSuite))
}

func (s *ReverseProxyTestSuite) TestRegisterReverseProxy() {
	s.Run("Happy Path", func() {
		backendResponse := "I am the backend"
		backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte(backendResponse))
		}))
		backendURL, _ := url.Parse(backend.URL)
		r := chi.NewRouter()
		frontends := make(configloader.Frontends)
		homeChild := make(map[string]configloader.Frontend)
		homeChild["start"] = configloader.Frontend{
			Url: backendURL.String(),
		}

		globalChild := make(map[string]configloader.Frontend)
		globalChild["footer"] = configloader.Frontend{
			Url: backendURL.String(),
		}
		frontends["global"] = globalChild
		frontends["home"] = homeChild
		RegisterReverseProxy(r, frontends)
		ts := httptest.NewServer(r)
		defer ts.Close()
		createdRoutes := r.Routes()

		routePattern := []string{}

		for _, one := range createdRoutes {
			routePattern = append(routePattern, one.Pattern)
			microUrlPattern := strings.ReplaceAll(one.Pattern, "*", "")
			handler := one.Handlers["GET"]
			req, _ := http.NewRequest("GET", fmt.Sprintf("%sassets/test.png", microUrlPattern), nil)
			res := newCloseNotifyingRecorder()
			handler.ServeHTTP(res, req)
			b, err := ioutil.ReadAll(res.Body)
			s.NoError(err)
			s.Equal(backendResponse, string(b))
		}

		s.ElementsMatch([]string{"/home.start/*", "/footer/*"}, routePattern)
	})
}
