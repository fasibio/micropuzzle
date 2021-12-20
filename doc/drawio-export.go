package main

import (
	"encoding/xml"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	files, err := LoadAllFilesWithExtension(".", "*.drawio")
	if err != nil {
		panic(err)
	}
	log.Println("Found", len(files), "files")
	for _, file := range files {
		name := strings.Replace(file.Name, ".drawio", "", 1)
		pages := getPagesOfDocument(string(file.Content))
		for i, page := range pages {
			cmd1 := exec.Command("drawio", "--export", "--format", "xml", "--page-index", strconv.Itoa(i), "--output", name+"_"+page.Name+".png", file.Name)
			cmd1.Stdout = os.Stdout
			err = cmd1.Run()
			if err != nil {
				panic(err)
			}
		}
	}
}

type DrawIoFile struct {
	Name    string
	Content []byte
}

func LoadAllFilesWithExtension(folder, pattern string) ([]DrawIoFile, error) {
	var matches []string
	var result []DrawIoFile
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	for _, match := range matches {
		content, err := loadFileFromDisk(match)
		if err != nil {
			panic(err)
		}
		result = append(result, DrawIoFile{Name: match, Content: content})
	}
	return result, err
}

type PageInformation struct {
	Name string `xml:"name,attr"`
	Id   string `xml:"id,attr"`
}

type XmlDigrams struct {
	Diagram []PageInformation `xml:"diagram"`
}

func getPagesOfDocument(content string) []PageInformation {
	var data XmlDigrams
	xml.Unmarshal([]byte(content), &data)
	return data.Diagram
}

func loadFileFromDisk(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}
