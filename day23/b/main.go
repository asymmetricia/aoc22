package main

import (
	"bytes"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/asymmetricia/aoc22/aoc"
	"github.com/asymmetricia/aoc22/coord"
	"github.com/asymmetricia/aoc22/isovox"
	"github.com/sirupsen/logrus"
)

const video = true
const renderGif = true
const scaled = true

var log = logrus.StandardLogger()

type world struct {
	world         *coord.SparseWorld
	proposal      map[coord.Coord][]coord.Coord
	elvesLastMove map[coord.Coord]int
	elvesColor    map[coord.Coord]color.Color
	cons          []consideration
}

func (w world) Render() image.Image {
	iw := isovox.World{map[isovox.Coord]*isovox.Voxel{}}

	ages := maps.Values(w.elvesLastMove)
	slices.Sort(ages)

	for _, elf := range w.world.Find('#') {
		var col color.Color
		if scaled {
			col = aoc.TolScale(ages[0], ages[len(ages)-1], w.elvesLastMove[elf])
		} else {
			col = w.elvesColor[elf]
		}
		iw.Voxels[isovox.Coord{elf.X, elf.Y, 0}] = &isovox.Voxel{Color: col}
	}

	img := iw.Render(9)
	if !renderGif {
		return img
	}

	imgp := image.NewPaletted(img.Bounds(), append(aoc.TolVibrant, palette.WebSafe...))
	draw.FloydSteinberg.Draw(imgp, img.Bounds(), img, image.Pt(0, 0))
	return imgp
}

func (w world) Clone() world {
	ret := world{
		world:         &coord.SparseWorld{},
		proposal:      maps.Clone(w.proposal),
		elvesLastMove: maps.Clone(w.elvesLastMove),
		elvesColor:    maps.Clone(w.elvesColor),
		cons:          slices.Clone(w.cons),
	}
	maps.Copy(*ret.world, *w.world)
	return ret
}

type consideration struct {
	Neighbors []coord.Direction
	Direction coord.Direction
}

func (c consideration) consider(w *world, elves []coord.Coord) []coord.Coord {
	var ret []coord.Coord
elves:
	for _, elf := range elves {
		hasNeigh := false
		for _, neigh := range elf.Neighbors(true) {
			if w.world.At(neigh) == '#' {
				hasNeigh = true
				break
			}
		}

		if !hasNeigh {
			// elf does not need to move
			continue
		}

		for _, dir := range c.Neighbors {
			if w.world.At(elf.Move(dir)) == '#' {
				ret = append(ret, elf)
				continue elves
			}
		}
		pos := elf.Move(c.Direction)
		w.proposal[pos] = append(w.proposal[pos], elf)
	}
	return ret
}

// If no other Elves are in one of those eight positions, the Elf does not do
// anything during this round. Otherwise, the Elf looks in each of four
// directions in the following order and proposes moving one step in the first
// valid direction:
//
//    If there is no Elf in the N, NE, or NW adjacent positions, the Elf proposes moving north one step.
//    If there is no Elf in the S, SE, or SW adjacent positions, the Elf proposes moving south one step.
//    If there is no Elf in the W, NW, or SW adjacent positions, the Elf proposes moving west one step.
//    If there is no Elf in the E, NE, or SE adjacent positions, the Elf proposes moving east one step.

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	uniq := map[string]bool{}
	for _, line := range lines {
		uniq[line] = true
	}
	log.Printf("read %d %s lines (%d unique)", len(lines), name, len(uniq))

	w := &world{
		world:         coord.Load(lines, false).(*coord.SparseWorld),
		proposal:      map[coord.Coord][]coord.Coord{},
		elvesLastMove: map[coord.Coord]int{},
		elvesColor:    map[coord.Coord]color.Color{},
		cons: []consideration{
			{[]coord.Direction{coord.North, coord.NorthEast, coord.NorthWest}, coord.North},
			{[]coord.Direction{coord.South, coord.SouthEast, coord.SouthWest}, coord.South},
			{[]coord.Direction{coord.West, coord.NorthWest, coord.SouthWest}, coord.West},
			{[]coord.Direction{coord.East, coord.NorthEast, coord.SouthEast}, coord.East},
		}}

	for _, dot := range w.world.Find('.') {
		w.world.Set(dot, 0)
	}

	for _, elf := range w.world.Find('#') {
		cols := []color.RGBA{
			aoc.TolVibrantMagenta,
			aoc.TolVibrantCyan,
			aoc.TolVibrantBlue,
			aoc.TolVibrantTeal,
			aoc.TolVibrantOrange,
			aoc.TolVibrantRed,
		}
		w.elvesColor[elf] = cols[rand.Intn(len(cols))]
		w.elvesLastMove[elf] = 0
	}

	ec := len(w.world.Find('#'))

	last := time.Now()
	count := 0
	var snapshots []world
	round := func() bool {
		log := log.WithField("round", count+1)
		init := w.world.Find('#')
		if len(init) != ec {
			panic("elf count changed")
		}
		elves := init
		for _, cons := range w.cons {
			elves = cons.consider(w, elves)
		}
		w.cons[0], w.cons[1], w.cons[2], w.cons[3] = w.cons[1], w.cons[2], w.cons[3], w.cons[0]
		if len(init) == len(elves) {
			return false
		}

		moved := 0
		for to, from := range w.proposal {
			if len(from) != 1 {
				continue
			}
			if w.world.At(from[0]) != '#' {
				log.Fatalf("no elf at from pos %v", from[0])
			}
			w.elvesColor[to] = w.elvesColor[from[0]]
			w.elvesLastMove[to] = count
			delete(w.elvesLastMove, from[0])
			w.world.Set(from[0], 0)
			w.world.Set(to, '#')
			moved++
			//log.Printf("from %v to %v, %d elves", from[0], to, len(*w.world))
		}

		for time.Since(last) > time.Second {
			log.Printf("round %d, %d elves moved of %d proposed moves", count+1, moved, len(w.proposal))
			last = time.Now()
		}

		if name == "test" || count%5 == 0 {
			snapshots = append(snapshots, w.Clone())
		}

		if len(w.proposal) == 0 {
			return false
		}

		w.proposal = map[coord.Coord][]coord.Coord{}
		count++

		return true
	}

	for round() {
	}

	if scaled {
		cycles := 0
		done := false
		for !done {
			cycles++
			done = true
			min := aoc.Min(maps.Values(w.elvesLastMove)...)
			for i, elf := range w.elvesLastMove {
				if elf > min {
					done = false
					w.elvesLastMove[i]--
				}
			}

			if cycles%5 == 0 {
				snapshots = append(snapshots, w.Clone())
			}
			if time.Since(last) >= time.Second {
				log.Printf("cooling... %d", cycles)
				last = time.Now()
			}
		}
	}

	mu := &sync.Mutex{}
	rendering := int32(0)
	images := make([]image.Image, len(snapshots))
	wg := &sync.WaitGroup{}
	ncpu := runtime.NumCPU()
	for i, w := range snapshots {
		for atomic.LoadInt32(&rendering) > int32(ncpu*2) {
			time.Sleep(time.Millisecond)
		}
		atomic.AddInt32(&rendering, 1)
		wg.Add(1)
		go func(i int, w world) {
			img := w.Render()
			mu.Lock()
			images[i] = img
			if time.Since(last) >= time.Second {
				log.Printf("rendering... %d/%d", i, len(snapshots))
				last = time.Now()
			}
			atomic.AddInt32(&rendering, -1)
			mu.Unlock()
			wg.Done()
		}(i, w)
	}
	wg.Wait()

	images = append(images, w.Render())

	if video {
		enc, err := aoc.NewMP4Encoder("day23-b-"+name+".mp4", 60, log)
		if err != nil {
			log.Fatal(err)
		}

		rect := images[0].Bounds()
		for _, img := range images {
			ir := img.Bounds()
			if ir.Dx() > rect.Dx() {
				rect.Min.X, rect.Max.X = ir.Min.X, ir.Max.X
			}
			if ir.Dy() > rect.Dy() {
				rect.Min.Y, rect.Max.Y = ir.Min.Y, ir.Max.Y
			}
		}
		padded := image.NewRGBA(rect)

		perc := 0
		for i, img := range images {
			//y := rect.Dx() * img.Bounds().Dy() / img.Bounds().Dx()
			//resized := resize.Resize(uint(rect.Dx()), uint(y), img, resize.Bicubic)
			//diff := rect.Dy() - y
			diffX := rect.Dx() - img.Bounds().Dx()
			diffY := rect.Dy() - img.Bounds().Dy()
			draw.Draw(padded, padded.Bounds(), image.Black, image.Pt(0, 0), draw.Over)
			draw.Draw(
				padded,
				image.Rect(diffX/2, diffY/2, diffX/2+img.Bounds().Dx(), diffY/2+img.Bounds().Dy()),
				img,
				image.Pt(0, 0),
				draw.Over,
			)

			enc.Encode(padded)
			if p := (i + 1) * 10 / len(images); p > perc {
				log.Printf("Encoding %d%%...", p*10)
				perc = p
			}
		}

		enc.Close()
	}

	if renderGif {
		anim := &gif.GIF{}
		for _, img := range images {
			anim.Image = append(anim.Image, img.(*image.Paletted))
			anim.Delay = append(anim.Delay, 2)
			anim.Disposal = append(anim.Disposal, gif.DisposalNone)
		}

		anim.Image = append(anim.Image, w.Render().(*image.Paletted))
		anim.Delay = append(anim.Delay, 2)
		anim.Disposal = append(anim.Disposal, gif.DisposalNone)
		anim.Delay[len(anim.Delay)-2] = 300

		aoc.SaveGIF(anim, "day23-b-"+name+".gif", log)
	}

	aoc.RenderPng(w.Render(), "day23-b-"+name+".png")

	return count + 1
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

	input := aoc.Input(2022, 23)
	log.Printf("input solution: %d", solution("input", input))
}
