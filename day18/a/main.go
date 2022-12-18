package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc22/aoc"
)

var log = logrus.StandardLogger()

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	log.Printf("read %d %s lines", len(lines), name)

	type coord struct {
		x, y, z int
	}
	world := map[coord]bool{}
	for _, line := range lines {
		var x, y, z int
		fmt.Sscanf(line, "%d,%d,%d", &x, &y, &z)
		world[coord{x, y, z}] = true
	}

	var surfaces int
	for cube := range world {
		for _, dx := range []int{-1, 1} {
			if !world[coord{cube.x+dx, cube.y, cube.z}] {
				surfaces++
			}
		}
		for _, dy := range []int{-1, 1} {
			if !world[coord{cube.x, cube.y+dy, cube.z}] {
				surfaces++
			}
		}
		for _, dz := range []int{-1, 1} {
			if !world[coord{cube.x, cube.y, cube.z+dz}] {
				surfaces++
			}
		}
	}

	return surfaces
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

	input := aoc.Input(2022, 18)
	log.Printf("input solution: %d", solution("input", input))
}
