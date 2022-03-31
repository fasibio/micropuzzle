package configloader

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

type Frontends struct {
	Definitions map[string]map[string]Frontend `yaml:"definitions"`
	Pages       Pages                          `yaml:"pages"`
}

type Pages map[string]Page

func (p Pages) GetKeyType() string {
	res := ""
	for k := range p {
		res += fmt.Sprintf("'%s'|", k)
	}
	res = strings.TrimRight(res, "|")
	return res
}

func (p Pages) GetPageByUrl(url string) *Page {
	for _, v := range p {
		if v.Url == url {
			return &v
		}
	}
	return nil
}

type Page struct {
	Url       string            `yaml:"url" json:"url,omitempty"`
	Title     string            `yaml:"title" json:"title,omitempty"`
	Fragments map[string]string `yaml:"fragments" json:"fragments,omitempty"`
}

func (p Page) GetFragmentByName(name string) string {
	return p.Fragments[name]
}

type Frontend struct {
	Url            string `yaml:"url"`
	GlobalOverride string `yaml:"globalOverride"`
}

// Find out url from configuration yaml by point seperated name (f.e. "startpage.content")
// See frontends.yaml to understand syntax
func (f Frontends) GetUrlByFrontendName(name string) string {
	val := strings.Split(name, ".")
	group := "global"
	if len(val) > 1 {
		group = val[0]
	}
	return f.Definitions[group][val[len(val)-1]].Url
}

func (f Frontends) GetKeyList() []string {
	var keys []string
	for k, v := range f.Definitions {
		for frontend := range v {
			keys = append(keys, k+"."+frontend)
		}
	}
	sort.Strings(keys)
	return keys
}

func (f Frontends) GetPagesList() map[string]Page {
	res := make(map[string]Page)

	globals := f.Pages["global"]
	for k, v := range f.Pages {
		if k != "global" {
			for k1, v1 := range globals.Fragments {
				if _, ok := v.Fragments[k1]; !ok {
					v.Fragments[k1] = v1
				}
			}
			res[k] = v
		}
	}
	return res
}

func LoadFrontends(frontendsPath string) (*Frontends, error) {
	frontendsBody, err := ioutil.ReadFile(frontendsPath)
	if err != nil {
		return nil, err
	}
	var frontends Frontends
	err = yaml.Unmarshal(frontendsBody, &frontends)
	if err != nil {
		return nil, err
	}
	return &frontends, nil
}
