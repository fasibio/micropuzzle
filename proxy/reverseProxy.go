package proxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/fasibio/micropuzzle/configloader"
	"github.com/fasibio/micropuzzle/logger"
	"github.com/fasibio/micropuzzle/resultmanipulation"
	"github.com/go-chi/chi"
)

func RegisterReverseProxy(r chi.Router, frontends configloader.Frontends) {
	for key, one := range frontends {
		for oneK, oneV := range one {
			prefix := ""
			if key != "global" {
				prefix = key + "."
			}
			err := registerMicrofrontendProxy(r, prefix+oneK, oneV)
			if err != nil {
				logger.Get().Warnw(fmt.Sprintf("Error by setting Reverseproxy for destination %s", prefix+oneK), "error", err)
			}
		}
	}
}

func registerMicrofrontendProxy(r chi.Router, name string, frontend configloader.Frontend) error {
	url, err := url.Parse(frontend.Url)
	if err != nil {
		return err
	}
	logger.Get().Infof("Register endpoint /%s/* for frontend %s", name, name)
	r.HandleFunc(fmt.Sprintf("/%s/*", name), func(w http.ResponseWriter, r *http.Request) {
		if frontend.StripUrlPrefix {
			path := strings.Replace(r.URL.String(), "/"+name, "", 1)
			r.URL, err = url.Parse(path)
		}
		p := httputil.NewSingleHostReverseProxy(url)
		p.ModifyResponse = rewriteBodyHandler("/" + name)
		p.ServeHTTP(w, r)
	})
	return nil
}

func rewriteBodyHandler(prefix string) func(*http.Response) error {
	return func(resp *http.Response) (err error) {
		b, err := ioutil.ReadAll(resp.Body) //Read html
		if err != nil {
			return err
		}
		log.Println("rewriteBodyHandler", resp.Request.URL.String(), resp.StatusCode)
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
