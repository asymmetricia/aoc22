package main

import (
	"bytes"
	"os"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc22/aoc"
)

var log = logrus.StandardLogger()

func ToSnafu(i int) string {
	if i == 0 {
		return ""
	}

	switch i % 5 {
	case 0:
		return ToSnafu(i/5) + "0"
	case 1:
		return ToSnafu(i/5) + "1"
	case 2:
		return ToSnafu(i/5) + "2"
	case 3:
		return ToSnafu(i/5+1) + "="
	default:
		return ToSnafu(i/5+1) + "-"
	}
}

func ParseSnafu(s string) int {
	var ret int
	for i, c := range s {
		n := map[rune]int{
			'2': 2,
			'1': 1,
			'0': 0,
			'-': -1,
			'=': -2,
		}[c]
		for j := 0; j < len(s)-i-1; j++ {
			n *= 5
		}
		ret += n
	}
	return ret
}

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	uniq := map[string]bool{}
	for _, line := range lines {
		uniq[line] = true
	}
	log.Printf("read %d %s lines (%d unique)", len(lines), name, len(uniq))

	sum := 0
	for _, line := range lines {
		sum += ParseSnafu(line)
	}

	log.Print(ToSnafu(sum))
	return -1
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})
	test, err := os.ReadFile("test")
	if err == nil {
		log.Printf("test solution: %d", solution("test", test))
	} else {
		log.Warningf("no test data present")
	}

	input := aoc.Input(2022, 25)
	log.Printf("input solution: %d", solution("input", input))
}
