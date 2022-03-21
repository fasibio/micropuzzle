package proxy

import (
	"compress/gzip"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/pquerna/cachecontrol"
)

var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

type CacheInformation struct {
	Expires time.Duration
	Header  http.Header
}

type proxy struct {
}

func NewProxy() *proxy {
	return &proxy{}
}

func (p *proxy) delHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
}

func (p *proxy) copyHeader(dst, src *http.Header) {
	for k, vv := range *src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func (p *proxy) appendHostToXForwardHeader(header http.Header, host string) {
	// If we aren't the first proxy retain prior
	// X-Forwarded-For information as a comma+space
	// separated list and fold multiple headers into one.
	if prior, ok := header["X-Forwarded-For"]; ok {
		host = strings.Join(prior, ", ") + ", " + host
	}
	header.Set("X-Forwarded-For", host)
}

func (p *proxy) Get(url string, header http.Header, remoteAddr string) ([]byte, CacheInformation, error) {
	client := &http.Client{}
	p.delHopHeaders(header)

	if clientIP, _, err := net.SplitHostPort(remoteAddr); err == nil {
		p.appendHostToXForwardHeader(header, clientIP)
	}
	req1, _ := http.NewRequest("GET", url, nil)
	p.copyHeader(&req1.Header, &header)
	resp, err := client.Do(req1)

	if err != nil {
		return nil, CacheInformation{Expires: time.Duration(0)}, err
	}
	defer resp.Body.Close()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}
	_, expires, _ := cachecontrol.CachableResponse(req1, resp, cachecontrol.Options{})
	diff := time.Until(expires)
	if diff < 0 {
		diff = time.Duration(0)

	}
	content, err := io.ReadAll(reader)
	return content, CacheInformation{Expires: diff, Header: resp.Header}, err
}
