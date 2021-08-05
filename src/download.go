package main

import (
	"compress/gzip"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	flag.Parse()
	id := flag.Arg(0)

	if id == "" {
		byAll()
	} else {
		byPlayerID(id)
	}
}

func byAll() {
	files := getLogFiles()

	for _, file := range files {
		html := getLogsHTML(file)
		ids := getMJLogIDs(html)

		for _, id := range ids {
			downloadMJLog(id, "")
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func byPlayerID(playerID string) {
	html := getLogsHTMLByPlayerID(playerID)
	ids := getMJLogIDs(html)

	for _, id := range ids {
		downloadMJLog(id, playerID)
		time.Sleep(time.Millisecond * 500)
	}
}

func downloadMJLog(id string, playerID string) {
	var file string
	if playerID == "" {
		file = "../mjlogs/" + id + ".mjlog"
	} else {
		file = "../mjlogs/" + playerID + "/" + id + ".mjlog"
		os.MkdirAll(filepath.Dir(file), 0755)
	}

	_, err := os.Stat(file)
	if !os.IsNotExist(err) {
		log.Printf("%s is exists", file)
		return
	}

	url := "https://tenhou.net/0/log/?" + id
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(file, body, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s is downloaded", file)
}

func getMJLogIDs(html *goquery.Document) []string {
	ids := []string{}

	html.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			r := regexp.MustCompile(`log=([\w\-]+)`)
			matches := r.FindStringSubmatch(href)
			if len(matches) > 0 {
				ids = append(ids, matches[1])
			}
		}
	})

	return ids
}

func getLogsHTML(file string) *goquery.Document {
	url := "https://tenhou.net/sc/raw/dat/" + file
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	reader, err := gzip.NewReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

func getLogsHTMLByPlayerID(playerID string) *goquery.Document {
	res, err := http.Get("https://tenhou.net/0/log/find.cgi?un=" + playerID)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

func getLogFiles() []string {
	res, err := http.Get("https://tenhou.net/sc/raw/list.cgi")
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	r := regexp.MustCompile(`\w+\.html\.gz`)
	files := []string{}
	for _, matches := range r.FindAllStringSubmatch(string(body), -1) {
		files = append(files, matches[0])
	}

	return files
}
