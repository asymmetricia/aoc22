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
	uniq := map[string]bool{}
	for _, line := range lines {
		uniq[line] = true
	}
	log.Printf("read %d %s lines (%d unique)", len(lines), name, len(uniq))

	values := map[string]int{}
	for _, line := range lines {
		var monkey string
		var value int
		_, err := fmt.Sscanf(line, "%4s: %d", &monkey, &value)
		if err == nil {
			log.Printf("found %s = %d", monkey, value)
			values[monkey] = value
		}
	}
	changed := true
	for changed {
		changed = false
		for _, line := range lines {
			var monkey, op1, op2 string
			var op rune
			_, err := fmt.Sscanf(line, "%4s: %4s %c %4s", &monkey, &op1, &op, &op2)
			if err != nil {
				continue
			}
			if _, ok := values[monkey]; ok {
				continue
			}
			v1, ok1 := values[op1]
			v2, ok2 := values[op2]
			if !ok1 || !ok2 {
				continue
			}
			switch op {
			case '+':
				values[monkey] = v1+v2
			case '-':
				values[monkey] = v1-v2
			case '*':
				values[monkey] = v1*v2
			case '/':
				values[monkey] = v1/v2
			default:
				panic(string(op))
			}
			log.Printf("concluded %s => %d", monkey, values[monkey])
			changed = true
		}
	}

	if name == "test" && values["root"] != 152 {
		panic("nope")
	}
	return values["root"]
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

	input := aoc.Input(2022, 21)
	log.Printf("input solution: %d", solution("input", input))
}
