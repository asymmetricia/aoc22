package main

import (
	"bytes"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc22/aoc"
	"github.com/asymmetricia/aoc22/canvas"
)

var log = logrus.StandardLogger()

type Directory struct {
	Files       map[string]int
	Directories map[string]*Directory
}

func (d *Directory) DirectoryNames() []string {
	if d == nil || d.Directories == nil {
		return nil
	}
	var ret []string
	for dir := range d.Directories {
		ret = append(ret, dir)
	}
	sort.Strings(ret)
	return ret
}

func (d *Directory) FileNames() []string {
	if d == nil || d.Files == nil {
		return nil
	}
	var ret []string
	for file := range d.Files {
		ret = append(ret, file)
	}
	sort.Strings(ret)
	return ret
}

func parseFrame(line string, dir string, root *Directory) *canvas.Canvas {
	ret := &canvas.Canvas{}
	x := 0
	y := 3
	col := func(dir *Directory, name string) (width int) {
		var body string
		for _, dirName := range dir.DirectoryNames() {
			dirName += "/"
			body += dirName + "\n"
			if len(dirName) > width {
				width = len(dirName)
			}
		}
		for _, file := range dir.FileNames() {
			body += file + "\n"
			if len(file) >= width {
				width = len(file)
			}
		}
		canvas.TextBox{
			Top:        y,
			Left:       x,
			Title:      []rune(name),
			Body:       []rune(strings.TrimSpace(body)),
			BodyColor:  aoc.TolVibrantBlue,
			TitleColor: aoc.TolVibrantMagenta,
			FrameColor: aoc.TolVibrantCyan,
		}.On(ret)
		return width
	}

	width := col(root, "/")
	path := strings.Split(strings.TrimLeft(dir, "/"), "/")
	cursor := root
	for len(path) > 0 {
		if path[0] == "" {
			break
		}
		dir, ok := cursor.Directories[path[0]]
		if !ok {
			break
		}
		for yy, dirname := range cursor.DirectoryNames() {
			if dirname == path[0] {
				y += yy + 1
				ret.PrintAt(x+1, y, dirname+"/", aoc.TolVibrantMagenta)
				break
			}
		}
		x += width
		width = col(dir, path[0])
		cursor = dir
		path = path[1:]
	}

	canvas.TextBox{
		Top:        0,
		Left:       0,
		Title:      []rune("parsing..."),
		Body:       []rune(line),
		BodyColor:  aoc.TolVibrantOrange,
		TitleColor: aoc.TolVibrantTeal,
		FrameColor: aoc.TolVibrantRed,
	}.On(ret)
	return ret
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

	var frames []*canvas.Canvas
	for _, line := range lines {
		log.Print(line)
		fields := strings.Fields(line)
		if fields[0] == "$" && fields[1] == "cd" {
			if fields[2] == "/" {
				dir = ""
			} else if fields[2] == ".." {
				path := strings.Split(dir, "/")
				dir = strings.Join(path[:len(path)-1], "/")
			} else {
				dir += "/" + fields[2]
			}
		} else if fields[0] == "$" && fields[1] == "ls" {
		} else if fields[0] == "dir" {
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

		frames = append(frames, parseFrame(line, dir, root))
	}

	canvas.RenderGif(
		frames,
		map[int]float32{10: 30, 20: 20, 30: 10, 100: 3, math.MaxInt: 1.0 / 2},
		"day07b-"+name+".gif",
		log,
	)

	const (
		total        = 70000000
		neededUnused = 30000000
	)
	used := root.Size()
	bestSize := math.MaxInt
	for dir := range root.WithSizeAbove("", neededUnused-(total-used)) {
		if root.Dir(dir).Size() < bestSize {
			bestSize = root.Dir(dir).Size()
		}
	}
	return bestSize
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
