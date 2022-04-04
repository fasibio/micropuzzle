package proxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/fasibio/micropuzzle/configloader"
	"github.com/fasibio/micropuzzle/logger"
	"github.com/fasibio/micropuzzle/mimetypefinder"
	"github.com/fasibio/micropuzzle/resultmanipulation"
)

type HttpRegistration interface {
	HandleFunc(pattern string, handlerFn http.HandlerFunc)
}

func RegisterReverseProxy(r HttpRegistration, frontends *configloader.Configuration) {

	//TODO also add result by cache control to redis cache

	for key, one := range frontends.Definitions {
		for oneK, oneV := range one {
			prefix := ""
			prefix = key + "."
			err := registerMicrofrontendProxy(r, prefix+oneK, oneV)
			if err != nil {
				logger.Get().Warnw(fmt.Sprintf("Error by setting Reverseproxy for destination %s", prefix+oneK), "error", err)
			}
		}
	}
}

func registerMicrofrontendProxy(r HttpRegistration, name string, frontend configloader.Definition) error {
	url, err := url.Parse(frontend.Url)
	if err != nil {
		return err
	}
	logger.Get().Infow("Register endpoint for reverseproxy", "path", fmt.Sprintf("/%s/*", name), "frontend", name, "frontend_url", url)
	handler := func(w http.ResponseWriter, r *http.Request) {
		path := strings.Replace(r.URL.String(), "/"+name, "", 1)
		r.URL, err = url.Parse(path)
		p := httputil.NewSingleHostReverseProxy(url)
		p.ModifyResponse = rewriteBodyHandler("/" + name)
		p.ServeHTTP(w, r)
	}
	if frontend.GlobalOverride != "" {
		logger.Get().Infow("Register GlobalOverride endpoint for frontend", "path", fmt.Sprintf("%s/*", frontend.GlobalOverride), "frontend", name, "frontend_url", url)
		r.HandleFunc(fmt.Sprintf("%s/*", frontend.GlobalOverride), handler)
	}
	r.HandleFunc(fmt.Sprintf("/%s/*", name), handler)
	return nil
}

func isContentTypeManipulable(contentType string) bool {
	return strings.Contains(contentType, "text/html") || strings.Contains(contentType, "text/css") || strings.Contains(contentType, "application/javascript") || contentType == ""
}

func rewriteBodyHandler(prefix string) func(*http.Response) error {
	return func(resp *http.Response) (err error) {
		if !isContentTypeManipulable(resp.Header.Get("Content-Type")) {
			return nil
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		err = resp.Body.Close()
		if err != nil {
			return err
		}
		var res string
		res = resultmanipulation.ChangePathOfRessourcesCss(string(b), prefix)
		res = resultmanipulation.ChangePathOfRessourcesJsModule(res, prefix)
		if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
			res, err = resultmanipulation.ChangePathOfRessources(res, prefix)
			if err != nil {
				return err
			}
		}
		body := ioutil.NopCloser(bytes.NewReader([]byte(res)))
		resp.Body = body
		resp.ContentLength = int64(len(res))

		resp.Header.Set("Content-Length", strconv.Itoa(len(res)))
		if resp.Header.Get("Content-Type") == "" {
			resp.Header.Set("Content-Type", mimetypefinder.MimeTypeForFile(res))
		}
		return nil
	}
}
