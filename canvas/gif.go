package canvas

import (
	"image"
	"image/gif"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc22/aoc"
)

const DefaultFrameTiming = 2

// RenderGif creates a GIF from the stack of canvases. (TODO: motion blending?)
func RenderGif(canvases []*Canvas, filename string, log logrus.FieldLogger) {
	width, height := Bounds(canvases)
	anim := &gif.GIF{}

	type Frame struct {
		img    *image.Paletted
		timing int
	}

	frames := &sync.Map{}
	wg := &sync.WaitGroup{}

	for i, canvas := range canvases {
		var frameTiming float32 = DefaultFrameTiming
		if canvas.Timing != 0 {
			frameTiming = canvas.Timing
		}
		if frameTiming < DefaultFrameTiming {
			if i%(int(1.0/frameTiming/float32(DefaultFrameTiming))) != 0 && i != len(canvases)-1 {
				continue
			}
			frameTiming = DefaultFrameTiming
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

	anim.Image = append(anim.Image, canvases[len(canvases)-1].RenderRect(width, height))
	anim.Delay = append(anim.Delay, anim.Delay[len(anim.Delay)-1])
	anim.Disposal = append(anim.Disposal, gif.DisposalNone)

	anim.Delay[len(anim.Delay)-2] = 500

	aoc.SaveGIF(anim, filename, log)
	log.Printf("rendered %d frames", len(canvases))
}
