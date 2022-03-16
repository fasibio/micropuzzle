package resultmanipulation

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type LinkTag struct {
	Tag  string
	Type string
}

var changeLinkTags = []LinkTag{
	{Tag: "link", Type: "href"},
	{Tag: "script", Type: "src"},
	{Tag: "img", Type: "src"},
	{Tag: "a", Type: "href"},
	{Tag: "iframe", Type: "src"},
	{Tag: "embed", Type: "src"},
	{Tag: "source", Type: "src"},
	{Tag: "track", Type: "src"},
	{Tag: "video", Type: "src"},
	{Tag: "audio", Type: "src"},
}

func ChangePathOfRessources(html, prefix string) (string, error) {

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return html, err
	}

	for _, tag := range changeLinkTags {
		doc.Find(tag.Tag).Each(func(i int, s *goquery.Selection) {
			changeHtmlLink(s, prefix, tag.Type)
		})
	}
	doc.Find("form").Each(func(i int, s *goquery.Selection) {
		action, ok := s.Attr("action")
		if ok && !strings.HasPrefix(action, "http") {
			s.SetAttr("action", prefix+action)
		}
	})
	return doc.Html()
}

func ChangePathOfRessourcesCss(css, prefix string) string {
	r := regexp.MustCompile(`url\(/`)
	return string(r.ReplaceAll([]byte(css), []byte("url("+prefix+"/")))
}

func changeHtmlLink(s *goquery.Selection, prefix, tag string) {
	tagValue, ok := s.Attr(tag)
	if ok && !strings.HasPrefix(tagValue, "http") {
		s.SetAttr(tag, prefix+tagValue)
	}
}
