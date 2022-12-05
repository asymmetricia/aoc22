package aoc

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Before(haystack, needle string) string {
	pos := strings.Index(haystack, needle)
	if pos < 0 {
		return haystack
	}
	return haystack[:pos]
}

func After(haystack, needle string) string {
	pos := strings.Index(haystack, needle)
	if pos < 0 {
		return ""
	}
	return haystack[pos+len(needle):]
}

func Split2(haystack, needle string) (string, string) {
	pos := strings.Index(haystack, needle)
	if pos < 0 {
		return haystack, ""
	}
	return haystack[:pos], haystack[pos+len(needle):]
}

func Int(in string) int {
	i, err := strconv.Atoi(in)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "could not parse %q as string: %v", in, err)
		os.Exit(1)
	}
	return i
}
