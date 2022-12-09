package main

import (
	"bytes"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc22/coord"
)

var log = logrus.StandardLogger()

func moveTail(head, tail coord.Coord) coord.Coord {
	dx := head.X - tail.X
	dy := head.Y - tail.Y
	// horiz
	if dy == 0 {
		if dx == 2 {
			return tail.East()
		} else if dx == -2 {
			return tail.West()
		}
		return tail
	}
	// vert
	if dx == 0 {
		if dy == 2 {
			return tail.South()
		} else if dy == -2 {
			return tail.North()
		}
		return tail
	}

	if dx == 2 || dy == 2 || dx == -2 || dy == -2 {
		if dx > 0 {
			tail = tail.East()
		} else {
			tail = tail.West()
		}
		if dy > 0 {
			tail = tail.South()
		} else {
			tail = tail.North()
		}
		return tail
	}
	return tail
}
func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	log.Printf("read %d %s lines", len(lines), name)

	var rope [10]coord.Coord
	positions := map[coord.Coord]bool{
		rope[9]: true,
	}

	for _, line := range lines {
		fields := strings.Fields(line)
		c, err := strconv.Atoi(fields[1])
		if err != nil {
			panic(line)
		}

		for c > 0 {
			c--
			switch fields[0] {
			case "U":
				rope[0] = rope[0].North()
			case "D":
				rope[0] = rope[0].South()
			case "L":
				rope[0] = rope[0].West()
			case "R":
				rope[0] = rope[0].East()
			}
			for i := 1; i <= 9; i++ {
				rope[i] = moveTail(rope[i-1], rope[i])
			}
			positions[rope[9]] = true
		}
	}

	return len(positions)
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
