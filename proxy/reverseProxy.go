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
	"github.com/fasibio/micropuzzle/resultmanipulation"
)

type HttpRegistration interface {
	HandleFunc(pattern string, handlerFn http.HandlerFunc)
}

func RegisterReverseProxy(r HttpRegistration, frontends configloader.Frontends) {

	//TODO also add result by cache control to redis cache

	for key, one := range frontends {
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

func registerMicrofrontendProxy(r HttpRegistration, name string, frontend configloader.Frontend) error {
	url, err := url.Parse(frontend.Url)
	if err != nil {
		return err
	}
	logger.Get().Infof("Register endpoint /%s/* for frontend %s with url: %s", name, name, url)
	handler := func(w http.ResponseWriter, r *http.Request) {
		path := strings.Replace(r.URL.String(), "/"+name, "", 1)
		r.URL, err = url.Parse(path)
		p := httputil.NewSingleHostReverseProxy(url)
		p.ModifyResponse = rewriteBodyHandler("/" + name)
		p.ServeHTTP(w, r)
	}
	if frontend.GlobalOverride != "" {
		logger.Get().Infof("Register GlobalOverride endpoint %s/* for frontend %s with url: %s", frontend.GlobalOverride, name, url)
		r.HandleFunc(fmt.Sprintf("%s/*", frontend.GlobalOverride), handler)
	}
	r.HandleFunc(fmt.Sprintf("/%s/*", name), handler)
	return nil
}

func isContentTypeManipulable(contentType string) bool {
	return strings.Contains(contentType, "text/html") || strings.Contains(contentType, "text/css") || strings.Contains(contentType, "application/javascript")
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
		if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
			res, err = resultmanipulation.ChangePathOfRessources(string(b), prefix)
			if err != nil {
				return err
			}
		} else if strings.Contains(resp.Header.Get("Content-Type"), "text/css") {
			res = resultmanipulation.ChangePathOfRessourcesCss(string(b), prefix)
		} else if strings.Contains(resp.Header.Get("Content-Type"), "application/javascript") {
			res = resultmanipulation.ChangePathOfRessourcesJsModule(string(b), prefix)
			// res = string(b)
		} else {
			res = string(b)
		}
		body := ioutil.NopCloser(bytes.NewReader([]byte(res)))
		resp.Body = body
		resp.ContentLength = int64(len(b))
		resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
		return nil
	}
}
