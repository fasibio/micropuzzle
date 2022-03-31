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

var jsFromDoubleSemicolonRegexp *regexp.Regexp
var jsFromSingleSemicolonRegexp *regexp.Regexp
var jsImportSingleQuotaRegex *regexp.Regexp

var jsImportDoubleQuotaRegex *regexp.Regexp
var cssUrlNoSemicolonRegex *regexp.Regexp
var cssSingleSemicolonRegex *regexp.Regexp
var cssDoubleSemicolonRegex *regexp.Regexp

func init() {
	jsFromDoubleSemicolonRegexp = regexp.MustCompile(`from "/`)
	jsFromSingleSemicolonRegexp = regexp.MustCompile(`from '/`)
	jsImportSingleQuotaRegex = regexp.MustCompile(`import '/`)
	jsImportDoubleQuotaRegex = regexp.MustCompile(`import "/`)
	cssUrlNoSemicolonRegex = regexp.MustCompile(`url\(/`)
	cssSingleSemicolonRegex = regexp.MustCompile(`url\('/`)
	cssDoubleSemicolonRegex = regexp.MustCompile(`url\(\"/`)
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

func ChangePathOfRessourcesJsModule(js, prefix string) (res string) {
	res = js
	res = string(jsFromSingleSemicolonRegexp.ReplaceAll([]byte(res), []byte(`from '`+prefix+`/`)))
	res = string(jsFromDoubleSemicolonRegexp.ReplaceAll([]byte(res), []byte(`from "`+prefix+`/`)))
	res = string(jsImportSingleQuotaRegex.ReplaceAll([]byte(res), []byte(`import '`+prefix+`/`)))
	res = string(jsImportDoubleQuotaRegex.ReplaceAll([]byte(res), []byte(`import "`+prefix+`/`)))
	return
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
	doc.Find("script[type='module']:not([src])").Each(func(i int, s *goquery.Selection) {
		s.SetHtml(ChangePathOfRessourcesJsModule(s.Text(), `/`+prefix))
	})
	doc.Find("form").Each(func(i int, s *goquery.Selection) {
		action, ok := s.Attr("action")
		if ok && !strings.HasPrefix(action, "http") {
			s.SetAttr("action", prefix+action)
		}
	})
	return doc.Html()
}

func ChangePathOfRessourcesCss(css, prefix string) (res string) {
	res = css
	res = string(cssUrlNoSemicolonRegex.ReplaceAll([]byte(res), []byte("url("+prefix+"/")))
	res = string(cssSingleSemicolonRegex.ReplaceAll([]byte(res), []byte("url('"+prefix+"/")))
	res = string(cssDoubleSemicolonRegex.ReplaceAll([]byte(res), []byte("url(\""+prefix+"/")))
	return
}

func changeHtmlLink(s *goquery.Selection, prefix, tag string) {
	tagValue, ok := s.Attr(tag)
	if ok && !strings.HasPrefix(tagValue, "http") {
		s.SetAttr(tag, prefix+tagValue)
	}
}
