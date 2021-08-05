package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/thoas/go-funk"
)

const (
	GAME_TYPE_HUMAN = 0x0001
	GAME_TYPE_THREE = 0x0010
)

const (
	YAKU_RIICH = "1"
)

func main() {
	files := getFiles()
	riichCount := 0
	riichSuccessCount := 0

	for _, file := range files {
		log := getLogText(file)

		riichCount += getRiichCount(log)
		riichSuccessCount += getRiichSuccessCount(log)
	}

	fmt.Println(riichCount, riichSuccessCount)
}

func getRiichCount(log string) int {
	// リーチ時にロンされなかった場合のみ`step="2"`が記録される
	r := regexp.MustCompile(`REACH[^>]*?step="2"/>`)
	matches := r.FindAllString(log, -1)
	return len(matches)
}

func getRiichSuccessCount(log string) int {
	count := 0
	r := regexp.MustCompile(`<AGARI[^>]*? yaku="([\d,]+)"`)
	for _, matches := range r.FindAllStringSubmatch(log, -1) {
		if len(matches) > 0 {
			// [役A,役Aの飜数,役B,役Bの飜数,...]
			for i, yaku := range strings.Split(matches[1], ",") {
				if i%2 == 0 && yaku == YAKU_RIICH {
					count += 1
				}
			}
		}
	}

	return count
}

func getLogText(file string) string {
	body, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}

func getFiles() []string {
	files := []string{}

	files, err := filepath.Glob("../mjlogs/*.mjlog")
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
