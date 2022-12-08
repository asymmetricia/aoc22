package main

import (
	"bytes"
	"os"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	log.Printf("read %d %s lines", len(lines), name)

	trees := [][]int{}
	for y, row := range lines {
		trees = append(trees, make([]int, len(row)))
		for x, height := range row {
			trees[y][x] = int(height - '0')
		}
	}

	visible := map[int]map[int]bool{}

	for y, row := range trees {
		visible[y] = map[int]bool{}
		visible[y][0] = true
		last := row[0]
		for x, height := range row {
			if height > last {
				last = height
				visible[y][x] = true
			}
		}

		visible[y][len(row)-1] = true
		last = row[len(row)-1]
		for x := len(row) - 1; x >= 0; x-- {
			if row[x] > last {
				last = row[x]
				visible[y][x] = true
			}
		}
	}

	for x := 0; x < len(trees[0]); x++ {
		visible[0][x] = true
		last := trees[0][x]
		for y, row := range trees {
			if row[x] > last {
				last = row[x]
				visible[y][x] = true
			}
		}

		visible[len(trees)-1][x] = true
		last = trees[len(trees)-1][x]
		for y := len(trees) - 1; y >= 0; y-- {
			if trees[y][x] > last {
				last = trees[y][x]
				visible[y][x] = true
			}
		}
	}

	count := 0
	for _, row := range visible {
		for _, tree := range row {
			if tree {
				count++
			}
		}
	}

	return count
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
