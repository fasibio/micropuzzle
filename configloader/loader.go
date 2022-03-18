package configloader

import (
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type Frontends map[string]map[string]Frontend

type Frontend struct {
	Url string `yaml:"url"`
}

// Find out url from configuration yaml by point seperated name (f.e. "startpage.content")
// See frontends.yaml to understand syntax
func (f Frontends) GetUrlByFrontendName(name string) string {
	val := strings.Split(name, ".")
	group := "global"
	if len(val) > 1 {
		group = val[0]
	}
	return f[group][val[len(val)-1]].Url
}

func (f Frontends) GetKeyList() []string {
	var keys []string
	for k, v := range f {
		for frontend := range v {
			keys = append(keys, k+"."+frontend)
		}
	}
	return keys
}

func LoadFrontends(frontendsPath string) (Frontends, error) {
	frontendsBody, err := ioutil.ReadFile(frontendsPath)
	if err != nil {
		return nil, err
	}
	var frontends Frontends
	err = yaml.Unmarshal(frontendsBody, &frontends)
	if err != nil {
		return nil, err
	}
	return frontends, nil
}
