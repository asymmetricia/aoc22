package canvas

import (
	"image"
	"image/gif"
	"math"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc22/aoc"
)

func timing(n int, csPerFrame map[int]float32) float32 {
	// the best timing is the timing T such that there is no timing less than T and greater than N.
	best := math.MaxInt
	for t := range csPerFrame {
		if t > n && t < best {
			best = t
		}
	}

	v, ok := csPerFrame[best]
	if !ok {
		return 3
	}

	return v
}

func RenderGif(canvases []*Canvas, csPerFrame map[int]float32, filename string, log logrus.FieldLogger) {
	width, height := Bounds(canvases)
	anim := &gif.GIF{}

	type Frame struct {
		img    *image.Paletted
		timing int
	}

	frames := &sync.Map{}
	wg := &sync.WaitGroup{}

	for i, canvas := range canvases {
		var frameTiming float32 = 3
		if csPerFrame != nil {
			frameTiming = timing(i, csPerFrame)
			if frameTiming < 1 {
				if i%(int(1.0/frameTiming)) != 0 && i != len(canvases)-1 {
					continue
				}
				frameTiming = 3
			}
		}

		wg.Add(1)
		go func(i int, canvas *Canvas, frameTiming float32) {
			defer wg.Done()
			if len(canvases) < 10 || i%(len(canvases)/10) == 0 {
				log.Printf("rendering frame %d", i)
			}
			frames.Store(i, Frame{canvas.RenderRect(width, height), int(frameTiming)})
		}(i, canvas, frameTiming)
	}

	wg.Wait()

	for i := range canvases {
		frame, ok := frames.Load(i)
		if !ok {
			continue
		}
		anim.Image = append(anim.Image, frame.(Frame).img)
		anim.Delay = append(anim.Delay, frame.(Frame).timing)
		anim.Disposal = append(anim.Disposal, gif.DisposalNone)
	}

	lastImage := anim.Image[len(anim.Image)-1]
	img := image.NewPaletted(lastImage.Rect, lastImage.Palette)
	copy(img.Pix, lastImage.Pix)

	anim.Image = append(anim.Image, lastImage)
	anim.Delay = append(anim.Delay, anim.Delay[len(anim.Delay)-1])
	anim.Disposal = append(anim.Disposal, gif.DisposalNone)

	anim.Delay[len(anim.Delay)-2] = 500

	aoc.SaveGIF(anim, filename, log)
}
