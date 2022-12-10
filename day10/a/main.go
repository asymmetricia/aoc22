package main

import (
	"bytes"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

type CPU struct {
	X     int
	Cycle int
	SS    int
}

func (c *CPU) Noop() {
	c.Cycle++
	c.CheckSS()
}

func (c *CPU) Addx(i int) {
	c.Noop()
	c.Cycle++
	c.CheckSS()
	c.X += i
}

func (c *CPU) CheckSS() {
	switch c.Cycle {
	case 20:
		fallthrough
	case 60:
		fallthrough
	case 100:
		fallthrough
	case 140:
		fallthrough
	case 180:
		fallthrough
	case 220:
		log.Printf("@cycle %d, x=%d, ss=%d", c.Cycle, c.X, c.Cycle*c.X)
		c.SS += c.Cycle * c.X
	}
}

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	log.Printf("read %d %s lines", len(lines), name)

	cpu := &CPU{X: 1}
	for _, line := range lines {
		parts := strings.Fields(line)
		switch parts[0] {
		case "noop":
			cpu.Noop()
		case "addx":
			immed, err := strconv.Atoi(parts[1])
			if err != nil {
				panic(err)
			}
			cpu.Addx(immed)
		}
	}

	return cpu.SS
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
