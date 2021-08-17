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
	"github.com/thoas/go-funk"

	"github.com/now-coding/mahjong-riich-success/domain"
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
		logs := getMJLogs(html)

		for _, l := range logs {
			downloadMJLog(l, "")
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func byPlayerID(playerID string) {
	html := getLogsHTMLByPlayerID(playerID)
	logs := getMJLogs(html)

	for _, l := range logs {
		downloadMJLog(l, playerID)
		time.Sleep(time.Millisecond * 500)
	}
}

func downloadMJLog(mjLog domain.MJLog, playerID string) {
	var file string
	if playerID == "" {
		file = "../mjlogs/" + mjLog.ID + ".mjlog"
	} else {
		file = "../mjlogs/" + playerID + "/" + mjLog.ID + "&tw=" + mjLog.MyPosition + ".mjlog"
		os.MkdirAll(filepath.Dir(file), 0755)
	}

	_, err := os.Stat(file)
	if !os.IsNotExist(err) {
		log.Printf("%s is exists", file)
		return
	}

	url := "https://tenhou.net/0/log/?" + mjLog.ID
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

func getMJLogs(html *goquery.Document) []domain.MJLog {
	logs := []domain.MJLog{}

	html.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			r := regexp.MustCompile(`log=([\w\-]+)(&tw=([\d]))?`)
			matches := r.FindStringSubmatch(href)
			if len(matches) > 0 {
				logs = append(logs, domain.MJLog{
					ID:         matches[1],
					MyPosition: matches[3],
				})
			}
		}
	})

	return funk.Reverse(logs).([]domain.MJLog)
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

	return funk.Reverse(files).([]string)
}
