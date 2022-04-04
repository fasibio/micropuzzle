package configloader

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type Configuration struct {
	Version     int                              `yaml:"version" validate:"required"`
	Definitions map[string]map[string]Definition `yaml:"definitions" validate:"required"`
	Pages       Pages                            `yaml:"pages" validate:"required"`
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
	Url       string            `yaml:"url" json:"url,omitempty" validate:"required"`
	Title     string            `yaml:"title" json:"title,omitempty"`
	Fragments map[string]string `yaml:"fragments" json:"fragments,omitempty" validate:"required"`
	Template  string            `yaml:"template" json:"-"`
}

func (p Page) GetFragmentByName(name string) string {
	return p.Fragments[name]
}

type Definition struct {
	Url            string `yaml:"url" validate:"required"`
	GlobalOverride string `yaml:"globalOverride"`
}

// Find out url from configuration yaml by point seperated name (f.e. "startpage.content")
// See frontends.yaml to understand syntax
func (f Configuration) GetUrlByFrontendName(name string) string {
	val := strings.Split(name, ".")
	group := "global"
	if len(val) > 1 {
		group = val[0]
	}
	return f.Definitions[group][val[len(val)-1]].Url
}

func (f Configuration) GetKeyList() []string {
	var keys []string
	for k, v := range f.Definitions {
		for frontend := range v {
			keys = append(keys, k+"."+frontend)
		}
	}
	sort.Strings(keys)
	return keys
}

func (f Configuration) GetPagesList() map[string]Page {
	res := make(map[string]Page)

	globals := f.Pages["global"]
	for k, v := range f.Pages {
		if k != "global" {
			if v.Template == "" {
				v.Template = globals.Template
			}
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

func LoadConfig(frontendsPath string) (*Configuration, error) {
	frontendsBody, err := ioutil.ReadFile(frontendsPath)
	if err != nil {
		return nil, err
	}
	var frontends Configuration
	err = yaml.Unmarshal(frontendsBody, &frontends)
	if err != nil {
		return nil, err
	}
	err = ValidateConfig(&frontends)
	if err != nil {
		return nil, err
	}
	return &frontends, nil
}

func ValidateConfig(frontends *Configuration) error {
	for _, v := range frontends.Definitions {
		for _, v1 := range v {
			if err := validate.Struct(v1); err != nil {
				return err
			}
		}
	}

	g, ok := frontends.Pages["global"]
	if ok {
		if err := validate.Var(g.Template, "required"); err != nil {
			return err
		}
	}
	for k, v := range frontends.GetPagesList() {
		if err := validate.Struct(v); err != nil {
			log.Println(k)
			return err
		}
	}
	return validate.Struct(*frontends)
}
