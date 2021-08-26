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

var (
	verbose bool
)

func main() {
	flag.BoolVar(&verbose, "v", false, "verbose tsumo or ron")
	flag.Parse()
	id := flag.Arg(0)

	files := getFiles(id)
	gamesCount := 0
	riichCount := 0
	riichTsumoSuccessCount := 0
	riichRonSuccessCount := 0

	for _, file := range files {
		log := getLog(file)

		gamesCount += 1
		riichCount += log.GetRiichCount()

		tsumo, ron := log.GetRiichSuccessCounts()
		riichTsumoSuccessCount += tsumo
		riichRonSuccessCount += ron
	}

	riichSuccessCount := riichTsumoSuccessCount + riichRonSuccessCount
	fmt.Printf("ゲーム数: %d\n", gamesCount)
	fmt.Printf("リーチ回数: %d\n", riichCount)

	if verbose {
		fmt.Printf("リーチ成功回数: %d (T:%d, R:%d)\n", riichSuccessCount, riichTsumoSuccessCount, riichRonSuccessCount)
		fmt.Printf(
			"リーチ成功率: %.2f (T:%.2f, R:%.2f)\n",
			float64(riichSuccessCount)/float64(riichCount)*100,
			float64(riichTsumoSuccessCount)/float64(riichCount)*100,
			float64(riichRonSuccessCount)/float64(riichCount)*100,
		)
	} else {
		fmt.Printf("リーチ成功回数: %d\n", riichSuccessCount)
		fmt.Printf("リーチ成功率: %.2f\n", float64(riichSuccessCount)/float64(riichCount)*100)
	}
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
