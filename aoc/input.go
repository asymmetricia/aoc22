package aoc

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var paths = []string{
	".",
	"..",
	"../..",
}

func Input(year int, day int) []byte {
	cacheFile := fmt.Sprintf(".input.%d.%d", year, day)
	for _, path := range paths {
		cache, err := os.ReadFile(filepath.Join(path, cacheFile))
		if err == nil {
			return cache
		}
	}

	var session []byte
	var path string
	for _, path = range paths {
		var err error
		session, err = os.ReadFile(filepath.Join(path, "aoc.session"))
		if err == nil {
			break
		}
	}
	if session == nil {
		log.Fatal("could not find any aoc.session")
	}
	session = bytes.TrimSpace(session)

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day),
		nil)
	if err != nil {
		panic(err)
	}
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: string(session),
	})

	req.Header.Set("user-agent", "tricia-adventofcode@cernu.us")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != 200 {
		panic(res.Status)
	}

	input, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(filepath.Join(path, cacheFile), input, 0644); err != nil {
		panic(err)
	}

	return input
}
