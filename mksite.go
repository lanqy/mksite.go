package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/daryl/fmatter"
	"github.com/mgutz/ansi"
	"github.com/russross/blackfriday"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

type data struct {
	Title       string
	Description string
	Author      string
	Created     string
}

// Item struct
type Item struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Created string `json:"created"`
}

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Xmlns   string   `xml:"xmlns,attr"`
	Title   string   `xml:"title"`
	Link    Link     `xml:"link"`
	Summary string   `xml:"summary"`
	Entrys  []entry  `xml:"entry"`
}

type entry struct {
	Title   string `xml:"title"`
	Link    Link   `xml:"link"`
	Summary string `xml:"summary"`
	Author  string `xml:"author"`
}

type Link struct {
	XMLName xml.Name `xml:link`
	Href    string   `xml:"href,attr"`
}

func main() {

	// config json file
	configFile := "config.json"

	// data json file
	dataJSON := "data.json"

	// atom.xml file
	atomFile := "atom.xml"

	jsonFile, err := os.Open(configFile)

	checkErr(err)

	defer jsonFile.Close()

	var result map[string]interface{}

	const htmlExt = ".html"
	const mdExt = ".md"
	const HOME = "index.html"

	// placeholder
	const TITLE = "$title"
	const DESCRIPTION = "$description"
	const BODY = "$body"
	const LINK = "$link"
	const NAME = "$name"
	const POST = "$post"
	const CREATED = "$created"
	const SITENAME = "$sitename"
	const DOMAIN = "https://lanqy.xyz"
	const XMLNS = "http://www.w3.org/2005/Atom"

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &result)
	sourceDir := result["sourceDir"].(string)
	fromDir := strings.Replace(sourceDir, "*", "", 1)
	files, err := filepath.Glob(sourceDir)

	checkErr(err)

	targetDir := result["targetDir"].(string)
	templateFile := result["templateFile"].(string)
	indexTemplateFile := result["indexTemplateFile"].(string)
	itemTemplateFile := result["itemTemplateFile"].(string)
	staticDir := result["staticDir"].(string)
	sitename := result["siteName"].(string)

	createFolder(targetDir)

	// date regexp pattern
	re := regexp.MustCompile("(\\d\\d\\d\\d)(/|-)(0?[1-9]|1[012])(/|-)(0?[1-9]|[12][0-9]|3[01])")
	var listArr []map[string]interface{}

	beginTime := time.Now()

	for _, file := range files {

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

			replacer := strings.NewReplacer(TITLE, string(d.Title), DESCRIPTION, string(d.Description), BODY, html)

			htmlString := replacer.Replace(string(template))

			htmlFileName := strings.TrimSuffix(_file, ext)

			fileName := "index" + htmlExt

			if string(d.Created) != "" && re.MatchString(d.Created) {
				dateSplit := strings.Split(replace(d.Created, "/", "-", 3), "-")
				path := strings.Join(dateSplit, "/")
				fullPath = targetDir + "/" + path
			} else {
				fullPath = targetDir
			}

			fullPath = fullPath + "/" + htmlFileName

			fmt.Println(fullPath)

			createFolder(fullPath)

			var item = make(map[string]interface{})
			item["title"] = string(d.Title)
			item["created"] = string(d.Created)
			item["link"] = replace(string(fullPath), targetDir, "", 1)

			listArr = append(listArr, item)

			writeFile(fullPath+"/"+fileName, []byte(htmlString))

			if runtime.GOOS == "windows" {
				fmt.Printf("\nBuilding file from %s to "+fullPath+"/"+fileName+" done!\n", fromDir+htmlFileName+ext)
			} else {
				text := ansi.Color("\nBuilding file from "+fromDir+htmlFileName+ext+" to "+fullPath+"/"+fileName+" done!\n", "green+b")
				fmt.Println(text)
			}

		}
	}

	jsonbytes, _ := json.Marshal(listArr)

	writeFile(dataJSON, []byte(jsonbytes))

	dataFile, err := os.Open(dataJSON)

	checkErr(err)

	defer dataFile.Close()

	byteValues, _ := ioutil.ReadAll(dataFile)

	var results []Item

	json.Unmarshal([]byte(byteValues), &results)

	Copy(staticDir, targetDir) // copy static file into targetDir

	var buffer bytes.Buffer

	// sort json results
	sort.Slice(results, func(i, j int) bool {
		p1, err := strconv.Atoi(replace(results[i].Created, "-", "", 3)) // 2018-12-10 to 20181210
		checkErr(err)
		p2, err := strconv.Atoi(replace(results[j].Created, "-", "", 3))
		checkErr(err)
		return p1 > p2
	})

	feedrss := &Feed{Xmlns: XMLNS, Title: sitename, Link: Link{Href: DOMAIN}, Summary: sitename}

	for _, v := range results {
		items, _ := ioutil.ReadFile(itemTemplateFile)
		str := strings.NewReplacer(NAME, string(v.Title), LINK, string(v.Link), CREATED, string(v.Created))
		list := str.Replace(string(items))
		feedrss.Entrys = append(feedrss.Entrys, entry{string(v.Title), Link{Href: string(v.Link)}, string(v.Title), "lanqy"})
		buffer.WriteString(list)
	}

	op, err := xml.MarshalIndent(feedrss, "  ", "    ")

	checkErr(err)

	writeFile(targetDir+"/"+atomFile, []byte(xml.Header+string(op))) // create atom.xml

	index, _ := ioutil.ReadFile(indexTemplateFile)
	indexStr := strings.NewReplacer(POST, buffer.String(), SITENAME, sitename)
	indexContent := indexStr.Replace(string(index))

	writeFile(targetDir+"/"+HOME, []byte(indexContent))

	if runtime.GOOS == "windows" {
		fmt.Println("\nBuilding home file to " + targetDir + "/" + HOME)
		fmt.Printf("\nBuild in %s \n\n", time.Since(beginTime))
	} else {
		msg := ansi.Color("\nBuilding home file to "+targetDir+"/"+HOME, "green+b")
		fmt.Println(msg)
		runtime := ansi.Color("\nBuild in "+time.Since(beginTime).String(), "green+b")
		fmt.Println(runtime)
	}
}

// Write file
func writeFile(filePath string, fileData []byte) {
	ioutil.WriteFile(filePath, fileData, 0644)
}

// create Folder by filepath
func createFolder(filePath string) {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) { // file does not exist
			os.MkdirAll(filePath, os.ModePerm)
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

// Copy copies src to dest, doesn't matter if src is a directory or a file
func Copy(src, dest string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}
	return copy(src, dest, info)
}

// copy dispatches copy-funcs according to the mode.
// Because this "copy" could be called recursively,
// "info" MUST be given here, NOT nil.
func copy(src, dest string, info os.FileInfo) error {
	if info.Mode()&os.ModeSymlink != 0 {
		return lcopy(src, dest, info)
	}
	if info.IsDir() {
		return dcopy(src, dest, info)
	}
	return fcopy(src, dest, info)
}

// fcopy is for just a file,
// with considering existence of parent directory
// and file permission.
func fcopy(src, dest string, info os.FileInfo) error {

	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = os.Chmod(f.Name(), info.Mode()); err != nil {
		return err
	}

	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	_, err = io.Copy(f, s)
	return err
}

// dcopy is for a directory,
// with scanning contents inside the directory
// and pass everything to "copy" recursively.
func dcopy(srcdir, destdir string, info os.FileInfo) error {

	if err := os.MkdirAll(destdir, info.Mode()); err != nil {
		return err
	}

	contents, err := ioutil.ReadDir(srcdir)
	if err != nil {
		return err
	}

	for _, content := range contents {
		cs, cd := filepath.Join(srcdir, content.Name()), filepath.Join(destdir, content.Name())
		if err := copy(cs, cd, content); err != nil {
			// If any error, exit immediately
			return err
		}
	}
	return nil
}

// lcopy is for a symlink,
// with just creating a new symlink by replicating src symlink.
func lcopy(src, dest string, info os.FileInfo) error {
	src, err := os.Readlink(src)
	if err != nil {
		return err
	}
	return os.Symlink(src, dest)
}
