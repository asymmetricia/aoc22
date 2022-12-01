package aoc

import (
	"bytes"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// Numbers1D loads and parses the input as a slice of integers
func Numbers1D() []int {
	input, err := ioutil.ReadFile("input")
	if err != nil {
		logrus.WithError(err).Fatal("could not parse input")
	}
	input = bytes.TrimSpace(input)
	lines := strings.Split(string(input), "\n")
	var nums []int
	for i, line := range lines {
		log := logrus.WithField("line", i).WithField("text", line)
		num, err := strconv.Atoi(line)
		if err != nil {
			log.WithError(err).Fatal("could not parse line")
		}
		nums = append(nums, num)
	}
	return nums
}

// Unique returns a new slice consisting of the objects in `in` in some order,
// attempting to avoid copies in the process.
func Unique(in []int) []int {
	if len(in) == 0 {
		return in
	}

	sort.Ints(in)
	cursor := -1
	for i := 0; i < len(in); i++ {
		if cursor == -1 || in[cursor] != in[i] {
			cursor++
			in[cursor] = in[i]
		}
	}
	in = in[:cursor+1]
	return in
}
