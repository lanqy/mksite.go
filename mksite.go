package main

import (
	"encoding/json"
	"fmt"
	"github.com/daryl/fmatter"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {

	configFile := "config.json"

	dataJSON := "data.json"

	jsonFile, err := os.Open(configFile)

	checkErr(err)

	type data struct {
		Title       string
		Description string
		Author      string
		Created     string
	}

	defer jsonFile.Close()

	var result map[string]interface{}

	var htmlExt = ".html"
	var mdExt = ".md"
	var titlePlaceholder = "$title"
	var descriptionPlaceholder = "$description"
	var bodyPlaceholder = "$body"

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &result)
	sourceDir := result["sourceDir"].(string)
	fromDir := strings.Replace(sourceDir, "*", "", 1)
	files, err := filepath.Glob(sourceDir)

	checkErr(err)

	targetDir := result["targetDir"].(string)
	templateFile := result["templateFile"].(string)

	createFolder(targetDir)

	re := regexp.MustCompile("(\\d\\d\\d\\d)(/|-)(0?[1-9]|1[012])(/|-)(0?[1-9]|[12][0-9]|3[01])") // date regexp

	var listArr []map[string]interface{}

	for i, file := range files {

		ext := filepath.Ext(file)
		_file := filepath.Base(file)
		if ext == mdExt { // markdown file
			var d data
			var fullPath = targetDir
			body, _ := ioutil.ReadFile(file)
			template, _ := ioutil.ReadFile(templateFile)
			content, err := fmatter.Parse([]byte(body), &d)

			checkErr(err)

			html := string(blackfriday.MarkdownCommon([]byte(content)))

			replaceTitle := replace(string(template), titlePlaceholder, string(d.Title), 1)
			replaceDes := replace(string(replaceTitle), descriptionPlaceholder, string(d.Description), 1)
			htmlString := replace(string(replaceDes), bodyPlaceholder, html, 1)
			htmlFileName := strings.TrimSuffix(_file, ext)

			fileName := htmlFileName + htmlExt

			if string(d.Created) != "" && re.MatchString(d.Created) {

				date := strings.Split(replace(d.Created, "/", "-", 3), "-")

				yearPath := targetDir + "/" + date[0]

				yearMonthPath := targetDir + "/" + date[0] + "/" + date[1]

				fullPath = targetDir + "/" + date[0] + "/" + date[1] + "/" + date[2]

				createFolder(yearPath)

				createFolder(yearMonthPath)

				createFolder(fullPath)
			}

			var item = make(map[string]interface{})
			item["title"] = string(d.Title)
			item["link"] = fullPath + "/" + fileName

			fmt.Println(item)
			fmt.Println(i)

			listArr = append(listArr, item)

			ioutil.WriteFile(fullPath+"/"+fileName, []byte(htmlString), 0644)

			fmt.Printf("\nBuilding file from %s to "+fullPath+"/"+fileName+" done!\n", fromDir+htmlFileName+ext)
		}
	}

	jsonbytes, _ := json.Marshal(listArr)

	ioutil.WriteFile(dataJSON, []byte(jsonbytes), 0644)
}

func createFolder(name string) {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) { // file does not exist
			os.Mkdir(name, os.ModePerm)
		}
	}
}

func replace(str string, hold string, s string, n int) string {
	return strings.Replace(str, hold, s, n)
}

func checkErr(err error) { // check error
	if err != nil {
		log.Fatal(err)
	}
}
