package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/thoas/go-funk"

	"github.com/now-coding/mahjong-riich-success/domain"
)

const (
	GAME_TYPE_HUMAN = 0x0001
	GAME_TYPE_THREE = 0x0010
)

func main() {
	flag.Parse()
	id := flag.Arg(0)

	files := getFiles(id)
	riichCount := 0
	riichSuccessCount := 0

	for _, file := range files {
		log := getLog(file)

		riichCount += log.GetRiichCount()
		riichSuccessCount += log.GetRiichSuccessCount()
	}

	fmt.Printf("リーチ回数: %d\n", riichCount)
	fmt.Printf("リーチ成功回数: %d\n", riichSuccessCount)
	fmt.Printf("リーチ成功率: %.2f\n", float64(riichSuccessCount)/float64(riichCount)*100)
}

func getLog(file string) domain.MJLog {
	body, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	filename := filepath.Base(file)
	r := regexp.MustCompile(`([\w\-]+)(&tw=([\d]))?`)
	matches := r.FindStringSubmatch(filename)

	return domain.MJLog{
		ID:         matches[1],
		MyPosition: matches[3],
		Body:       string(body),
	}
}

func getFiles(id string) []string {
	var files []string
	var err error

	if id == "" {
		files, err = filepath.Glob("../mjlogs/*.mjlog")
	} else {
		files, err = filepath.Glob("../mjlogs/" + id + "/*.mjlog")
	}
	if err != nil {
		log.Fatal(err)
	}

	return funk.FilterString(files, func(file string) bool {
		filename := filepath.Base(file)
		names := strings.Split(filename, "-")

		// 4麻対人戦のみに限定する
		// SEE: https://m77.hatenablog.com/entry/2017/05/21/214529
		types, err := strconv.ParseInt(names[1], 16, 32)
		if err != nil {
			log.Fatal(err)
		}

		return types&GAME_TYPE_HUMAN > 0 && types&GAME_TYPE_THREE == 0
	})
}
