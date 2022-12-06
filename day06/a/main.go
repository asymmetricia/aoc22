package main

import (
	"bytes"
	"os"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

func isMarker(s string) bool {
	if len(s) != 4 {
		panic(s)
	}
	if s[0] == s[1] || s[0] == s[2] || s[0] == s[3] ||
		s[1] == s[2] || s[1] == s[3] ||
		s[2] == s[3] {
		return false
	}
	return true
}

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	log.Printf("read %d %s lines", len(lines), name)

	for i := 0; i < len(input); i++ {
		if !isMarker(string(input[i : i+4])) {
			continue
		}
		log.Printf("%d, %s", i+4, input[i:i+4])
		break
	}
	return -1
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
