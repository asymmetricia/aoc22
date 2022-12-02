package main

import (
	"bytes"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

func outcome(game string) int {
	winner := map[byte]byte{
		'A': 'Y',
		'B': 'Z',
		'C': 'X',
	}
	loser := map[byte]byte{
		'A': 'Z',
		'B': 'X',
		'C': 'Y',
	}
	game = strings.Replace(game, " ", "", -1)
	switch game[1] {
	case 'X':
		game = string([]byte{game[0], loser[game[0]]})
	case 'Y':
		game = string([]byte{game[0], game[0] + 'X' - 'A'})
	case 'Z':
		game = string([]byte{game[0], winner[game[0]]})
	}
	switch game {
	case "AX":
		return 1 + 3
	case "AY":
		return 2 + 6
	case "AZ":
		return 3 + 0
	case "BX":
		return 1 + 0
	case "BY":
		return 2 + 3
	case "BZ":
		return 3 + 6
	case "CX":
		return 1 + 6
	case "CY":
		return 2 + 0
	case "CZ":
		return 3 + 3
	}
	return -1
}

func solution(input []byte) int {
	input = bytes.TrimSpace(input)
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	lines := strings.Split(strings.TrimSpace(string(input)), "\n")
	log.Printf("read %d input lines", len(lines))
	var score int
	for _, line := range lines {
		if o := outcome(line); o < 0 {
			panic(line)
		} else {
			score += o
		}
	}
	return score
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
