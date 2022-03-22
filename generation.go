package main

import (
	_ "embed"

	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/fasibio/micropuzzle/configloader"
	"github.com/urfave/cli/v2"
)

//go:embed mircopuzzle-clientlib/index.ts
var clientLibFile string

func (ru *runner) GenerateType(c *cli.Context) error {
	destination := c.String(CliGenDestination)
	sourceUrl := c.String(CliGenUrl)
	var keyList []string

	if sourceUrl == "" {
		frontends, err := configloader.LoadFrontends(c.String(CliMicrofrontends))
		if err != nil {
			return err
		}
		keyList = frontends.GetKeyList()
	} else {
		u, err := url.Parse(sourceUrl)
		if err != nil {
			return err
		}
		u.Path += path.Join(u.Path, "/frontends")
		res, err := http.Get(u.String())
		if err != nil {
			return err
		}
		var frontends FrontedsManagementResult
		json.NewDecoder(res.Body).Decode(&frontends)
		defer res.Body.Close()
		keyList = frontends.Frontends
	}

	destinationContent := getTypeScriptContent(keyList)
	if _, err := os.Stat(destination); os.IsNotExist(err) {
		err := os.Mkdir(destination, os.ModeDir|0755)
		if err != nil {
			return err
		}
	}
	os.Remove(fmt.Sprintf("%s/index.ts", destination))
	file, err := os.Create(fmt.Sprintf("%s/%s", destination, "index.ts"))
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(destinationContent)

	return nil
}

func getTypeScriptContent(keyList []string) string {
	destinationContent := `/**
 * Mircopuzzle AUTO-GENERATED CODE: PLEASE DO NOT MODIFY MANUALLY
 */
	
export enum MicropuzzleFrontends {
`
	for _, key := range keyList {
		destinationContent += fmt.Sprintf("\t%s=\"%s\",\n", strings.ToUpper(strings.Replace(key, ".", "_", 1)), key)
	}
	destinationContent += "}\n\n"
	destinationContent += string(clientLibFile)
	return destinationContent
}
