package main

import (
	"bytes"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc22/aoc"
)

var log = logrus.StandardLogger()

func solution(name string, input []byte) int {
	input = bytes.TrimSpace(input)
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	lines := strings.Split(strings.TrimSpace(string(input)), "\n")
	log.Printf("read %d input lines", len(lines))

	var contains int

	for _, line := range lines {
		pairs := strings.Split(line, ",")
		s1, e1 := aoc.Split2(pairs[0], "-")
		s1i, e1i := aoc.Int(s1), aoc.Int(e1)
		s2, e2 := aoc.Split2(pairs[1], "-")
		s2i, e2i := aoc.Int(s2), aoc.Int(e2)

		if s1i <= s2i && e1i >= e2i ||
			s2i <= s1i && e2i >= e1i {
			contains++
		}
	}

	return contains
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

	input, err := os.ReadFile("input")
	if err != nil {
		log.WithError(err).Fatal("could not read input")
	}
	log.Printf("input solution: %d", solution("input", input))
}
