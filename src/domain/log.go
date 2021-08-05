package domain

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/thoas/go-funk"
)

type MJLog struct {
	ID         string
	MyPosition string
	Body       string
}

func (m MJLog) GetRiichCount() int {
	r := regexp.MustCompile(`<REACH.*?/>`)
	matches := r.FindAllString(m.Body, -1)

	return len(funk.Filter(matches, func(match string) bool {
		// リーチ時にロンされなかった場合のみ`step="2"`が記録される
		if !strings.Contains(match, `step="2"`) {
			return false
		}

		// MyPositionの指定がある場合、自分自身のリーチのみを集計
		if m.MyPosition != "" && !strings.Contains(match, fmt.Sprintf(`who="%s"`, m.MyPosition)) {
			return false
		}

		return true
	}).([]string))
}

func (m MJLog) GetRiichSuccessCount() int {
	count := 0

	r := regexp.MustCompile(`<AGARI.*?/>`)
	agaries := r.FindAllString(m.Body, -1)
	for _, agari := range agaries {
		// MyPositionの指定がある場合、自分自身のあがりのみを集計
		if m.MyPosition != "" && !strings.Contains(agari, fmt.Sprintf(`who="%s"`, m.MyPosition)) {
			continue
		}

		r := regexp.MustCompile(`yaku="([\d,]+)"`)
		matches := r.FindStringSubmatch(agari)
		if len(matches) > 0 {
			// [役A,役Aの飜数,役B,役Bの飜数,...]
			for i, yaku := range strings.Split(matches[1], ",") {
				if i%2 == 0 && yaku == YAKU_RIICH {
					count += 1
				}
				if i%2 == 0 && yaku == YAKU_W_RIICH {
					count += 1
				}
			}
		}
	}

	return count
}
