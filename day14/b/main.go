package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc22/aoc"
	"github.com/asymmetricia/aoc22/coord"
)

var log = logrus.StandardLogger()

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	log.Printf("read %d %s lines", len(lines), name)

	world := &coord.DenseWorld{}
	for _, line := range lines {
		fields := strings.Split(line, " -> ")
		var prev coord.Coord
		for i, field := range fields {
			c := coord.MustFromComma(field)
			if i == 0 {
				prev = c
				world.Set(prev, '#')
				continue
			}
			for prev.X < c.X {
				prev.X++
				world.Set(prev, '#')
			}
			for prev.X > c.X {
				prev.X--
				world.Set(prev, '#')
			}
			for prev.Y < c.Y {
				prev.Y++
				world.Set(prev, '#')
			}
			for prev.Y > c.Y {
				prev.Y--
				world.Set(prev, '#')
			}
		}
	}

	var sandStart = coord.MustFromComma("500,0")
	world.Set(sandStart, 'v')

	y := len(*world) + 1

	particles := 0
	pos := sandStart
	var worlds []*coord.DenseWorld
	for {
		shouldPath := pos != sandStart && particles%5 == 0 && particles < 100
		if pos.South().Y < y {
			if world.At(pos.South()) == 0 {
				if shouldPath {
					world.Set(pos, 0)
					world.Set(pos.South(), '+')
					worlds = append(worlds, world.Crop())
				}
				pos = pos.South()
				continue
			}
			if world.At(pos.SouthWest()) == 0 {
				if shouldPath {
					world.Set(pos, 0)
					world.Set(pos.SouthWest(), '+')
					worlds = append(worlds, world.Crop())
				}
				pos = pos.SouthWest()
				continue
			}
			if world.At(pos.SouthEast()) == 0 {
				if shouldPath {
					world.Set(pos, 0)
					world.Set(pos.SouthEast(), '+')
					worlds = append(worlds, world.Crop())
				}
				pos = pos.SouthEast()
				continue
			}
		}
		if world.At(pos) == 'v' {
			particles++
			break
		}
		world.Set(pos, '+')
		particles++
		pos = sandStart
		// show 10% of frames below particle 2k & 2% of frames above
		if particles < 2000 && particles%10 == 0 ||
			particles%50 == 0 ||
			name == "test" {
			worlds = append(worlds, world.Crop())
		}
	}

	log.Print(particles)

	last := worlds[len(worlds)-1]
	vpos := last.Find('v')[0]
	_, _, maxx, maxy := last.Rect()

	anim := &gif.GIF{
		Image:    make([]*image.Paletted, len(worlds)),
		Delay:    make([]int, len(worlds)),
		Disposal: make([]byte, len(worlds)),
	}
	const scale = 4
	wg := &sync.WaitGroup{}
	for i, world := range worlds {
		wg.Add(1)
		go func(i int, world *coord.DenseWorld) {
			defer wg.Done()
			img := image.NewPaletted(image.Rect(0, 0, maxx*scale, maxy*scale), aoc.TolVibrant)
			wvpos := world.Find('v')[0]
			adj := vpos.Minus(wvpos)
			world.Each(func(c coord.Coord) (stop bool) {
				view := c.Plus(adj)
				var col color.Color = aoc.TolVibrantGrey
				switch world.At(c) {
				case 0:
					col = color.Black
				case '+':
					col = aoc.TolVibrantOrange
				default:
					col = aoc.TolVibrantGrey
				}
				draw.Draw(
					img,
					image.Rect(view.X*scale, view.Y*scale, (view.X+1)*scale, (view.Y+1)*scale),
					image.NewUniform(col),
					image.Pt(0, 0),
					draw.Over)
				return false
			})

			aoc.Typeset(img, image.Point{}, strconv.Itoa(len(world.Find('+'))), aoc.TolVibrantTeal, aoc.TypesetOpts{Scale: scale})

			anim.Image[i] = img
			anim.Delay[i] = 2
			anim.Disposal[i] = gif.DisposalNone
		}(i, world)
	}

	wg.Wait()

	aoc.SaveGIF(anim, "day14-"+name+"-b.gif", log)

	return particles
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

	input := aoc.Input(2022, 14)
	log.Printf("input solution: %d", solution("input", input))
}
