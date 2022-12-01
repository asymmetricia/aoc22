package main

import (
	"bytes"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

func solution(input []byte) int {
	input = bytes.TrimSpace(input)
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	lines := strings.Split(strings.TrimSpace(string(input)), "\n")
	log.Printf("read %d input lines", len(lines))

	return -1
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
