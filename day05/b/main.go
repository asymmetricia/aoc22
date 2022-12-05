package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc22/aoc"
)

var log = logrus.StandardLogger()

func move(stacks [][]byte, n, from, to int) [][]byte {
	if n > len(stacks[from]) {
		log.Fatalf("cannot move %d crates from stack %d %v", n, from, stacks[from])
	}
	stacks[to] = append(stacks[to], stacks[from][len(stacks[from])-n:len(stacks[from])]...)
	stacks[from] = stacks[from][:len(stacks[from])-n]
	return stacks
}

func solution(name string, input []byte) string {
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	lines := strings.Split(string(input), "\n")
	log.Printf("read %d input lines", len(lines))

	stacks := make([][]byte, 9)
	state := 0
	for lineno, line := range lines {
		switch state {
		case 0:
			if line == "" {
				state++
				continue
			}
			for i := 0; i < len(line); i += 4 {
				if line[i] == '[' {
					stacks[i/4] = append([]byte{line[i+1]}, stacks[i/4]...)
				}
			}

		case 1:
			if line == "" {
				state++
				continue
			}
			re := regexp.MustCompile(`move (\d+) from (\d+) to (\d+)`)
			matches := re.FindStringSubmatch(line)
			if matches == nil {
				panic(strconv.Itoa(lineno) + ": " + line)
			}
			n, from, to := aoc.Int(matches[1]), aoc.Int(matches[2])-1, aoc.Int(matches[3])-1
			stacks = move(stacks, n, from, to)
		}
	}

	ans := ""
	for n, stack := range stacks {
		if len(stack) > 0 {
			ans += string(stack[len(stack)-1])
		}
		m := fmt.Sprintf("%d: ", n+1)
		for _, crate := range stack {
			m += string(crate) + " "
		}
		log.Print(m)
	}

	return ans
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})
	test, err := os.ReadFile("test")
	if err == nil {
		log.Printf("test solution: %s", solution("test", test))
	} else {
		log.Warningf("no test data present")
	}

	input, err := os.ReadFile("input")
	if err != nil {
		log.WithError(err).Fatal("could not read input")
	}
	log.Printf("input solution: %s", solution("input", input))
}
