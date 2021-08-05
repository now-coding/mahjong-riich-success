package main

import (
	"compress/gzip"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	files := getLogFiles()

	for _, file := range files {
		html := getLogsHTML(file)
		ids := getMJLogIDs(html)

		for _, id := range ids {
			downloadMJLog(id)
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func downloadMJLog(id string) {
	file := "../mjlogs/" + id + ".mjlog"

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

	ioutil.WriteFile(file, body, 0666)
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
