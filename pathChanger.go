package main

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ChangePathOfRessources(html, prefix string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		return html, err
	}
	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok && !strings.HasPrefix(href, "http") {
			s.SetAttr("href", prefix+href)
		}
	})
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		src, ok := s.Attr("src")
		if ok && !strings.HasPrefix(src, "http") {
			s.SetAttr("src", prefix+src)
		}
	})
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, ok := s.Attr("src")
		if ok && !strings.HasPrefix(src, "http") {
			s.SetAttr("src", prefix+src)
		}
	})
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok && !strings.HasPrefix(href, "http") {
			s.SetAttr("href", prefix+href)
		}
	})
	doc.Find("form").Each(func(i int, s *goquery.Selection) {
		action, ok := s.Attr("action")
		if ok && !strings.HasPrefix(action, "http") {
			s.SetAttr("action", prefix+action)
		}
	})
	doc.Find("iframe").Each(func(i int, s *goquery.Selection) {
		src, ok := s.Attr("src")
		if ok && !strings.HasPrefix(src, "http") {
			s.SetAttr("src", prefix+src)
		}
	})
	doc.Find("embed").Each(func(i int, s *goquery.Selection) {
		src, ok := s.Attr("src")
		if ok && !strings.HasPrefix(src, "http") {
			s.SetAttr("src", prefix+src)
		}
	})
	doc.Find("source").Each(func(i int, s *goquery.Selection) {
		src, ok := s.Attr("src")
		if ok && !strings.HasPrefix(src, "http") {
			s.SetAttr("src", prefix+src)
		}
	})
	doc.Find("track").Each(func(i int, s *goquery.Selection) {
		src, ok := s.Attr("src")
		if ok && !strings.HasPrefix(src, "http") {
			s.SetAttr("src", prefix+src)
		}
	})
	doc.Find("video").Each(func(i int, s *goquery.Selection) {
		src, ok := s.Attr("src")
		if ok && !strings.HasPrefix(src, "http") {
			s.SetAttr("src", prefix+src)
		}
	})
	doc.Find("audio").Each(func(i int, s *goquery.Selection) {
		src, ok := s.Attr("src")
		if ok && !strings.HasPrefix(src, "http") {
			s.SetAttr("src", prefix+src)
		}
	})

	return doc.Html()
}

func ChangePathOfRessourcesCss(css, prefix string) string {
	r := regexp.MustCompile(`url\(/`)
	return string(r.ReplaceAll([]byte(css), []byte("url("+prefix+"/")))
}
