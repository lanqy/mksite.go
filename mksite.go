package main

import (
	"encoding/json"
	"fmt"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	configFile := "config.json"
	jsonFile, err := os.Open(configFile)
	checkErr(err)

	defer jsonFile.Close()
	var result map[string]interface{}
	var htmlExt = ".html"
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &result)
	sourceDir := result["sourceDir"].(string)
	fromDir := strings.Replace(sourceDir, "*", "", 1)
	files, err := filepath.Glob(sourceDir)
	checkErr(err)
	targetDir := result["targetDir"].(string)
	templateFile := result["templateFile"].(string)
	if _, err := os.Stat(targetDir); err != nil {
		if os.IsNotExist(err) { // file does not exist
			os.Mkdir(targetDir, os.ModePerm)
		}
	}

	for _, file := range files {

		ext := filepath.Ext(file)
		_file := filepath.Base(file)
		if ext == ".md" { // markdown file
			body, _ := ioutil.ReadFile(file)
			template, _ := ioutil.ReadFile(templateFile)
			html := string(blackfriday.MarkdownCommon([]byte(body)))
			htmlString := strings.Replace(string(template), "$body", html, 1)
			htmlFileName := strings.TrimSuffix(_file, ext)
			ioutil.WriteFile(targetDir+"/"+htmlFileName+htmlExt, []byte(htmlString), 0644)
			fmt.Printf("build file from %s to "+targetDir+"/"+htmlFileName+htmlExt+" done!\n\n", fromDir+htmlFileName+ext)
		}
	}
}

func checkErr(err error) { // check error
	if err != nil {
		panic(err)
	}
}
