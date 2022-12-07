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

type Directory struct {
	Files       map[string]int
	Directories map[string]*Directory
}

func (d *Directory) Dir(path string) *Directory {
	if len(path) == 0 {
		return d
	}
	parts := strings.Split(strings.TrimLeft(path, "/"), "/")
	if dir, ok := d.Directories[parts[0]]; ok {
		return dir.Dir(strings.Join(parts[1:], "/"))
	}
	return nil
}

func (d *Directory) Size() int {
	size := 0
	for _, filesize := range d.Files {
		size += filesize
	}
	for _, subdir := range d.Directories {
		size += subdir.Size()
	}
	return size
}

func (d *Directory) WithSizeAbove(path string, n int) map[string]bool {
	if len(d.Directories) == 0 {
		return nil
	}

	ret := map[string]bool{}

	for dirName, dir := range d.Directories {
		if dir.Size() > n {
			ret[path+"/"+dirName] = true
		}
		for above := range dir.WithSizeAbove(path+"/"+dirName, n) {
			ret[above] = true
		}
	}

	return ret
}

func (d *Directory) WithSizeAtMost(path string, n int) map[string]bool {
	if len(d.Directories) == 0 {
		return nil
	}

	ret := map[string]bool{}

	for dirName, dir := range d.Directories {
		if dir.Size() <= n {
			ret[path+"/"+dirName] = true
		}
		for above := range dir.WithSizeAtMost(path+"/"+dirName, n) {
			ret[above] = true
		}
	}

	return ret
}

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	log.Printf("read %d %s lines", len(lines), name)

	var dir string
	root := &Directory{
		Files:       map[string]int{},
		Directories: map[string]*Directory{},
	}
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "$" {
			switch fields[1] {
			case "cd":
				if fields[2] == "/" {
					dir = ""
				} else if fields[2] == ".." {
					path := strings.Split(dir, "/")
					dir = strings.Join(path[:len(path)-1], "/")
				} else {
					dir += "/" + fields[2]
				}
			case "ls":
			default:
				panic(line)
			}
			continue
		}

		if fields[0] == "dir" {
			root.Dir(dir).Directories[fields[1]] = &Directory{
				Files:       map[string]int{},
				Directories: map[string]*Directory{},
			}
		} else {
			var err error
			root.Dir(dir).Files[fields[1]], err = strconv.Atoi(fields[0])
			if err != nil {
				panic(line)
			}
		}

	}

	log.Print(root.WithSizeAtMost("", 100000))

	ans := 0
	for dir := range root.WithSizeAtMost("", 100000) {
		ans += root.Dir(dir).Size()
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
