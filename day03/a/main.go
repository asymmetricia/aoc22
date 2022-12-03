package main

import (
	"bytes"
	"github.com/asymmetricia/aoc22/set"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

func prio(r rune) int {
	if r >= 'a' && r <= 'z' {
		return int(r - 'a' + 1)
	}
	return int(r - 'A' + 27)
}

func solution(input []byte) int {
	input = bytes.TrimSpace(input)
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	lines := strings.Split(strings.TrimSpace(string(input)), "\n")
	log.Printf("read %d input lines", len(lines))

	var sum = 0

	for _, line := range lines {
		first := line[:len(line)/2]
		second := line[len(line)/2:]

		fset := set.FromString(first)
		sset := set.FromString(second)
		for dupe := range fset.Intersect(sset) {
			sum += prio(dupe)
		}
	}

	return sum
}

func main() {
	test, err := os.ReadFile("test")
	if err == nil {
		log.Printf("test solution: %d", solution(test))
	} else {
		log.Warningf("no test data present")
	}

	input, err := os.ReadFile("input")
	if err != nil {
		log.WithError(err).Fatal("could not read input")
	}
	log.Printf("input solution: %d", solution(input))
}
