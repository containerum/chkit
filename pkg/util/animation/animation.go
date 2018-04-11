package animation

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

type FrameSource interface {
	Next() bool
	Frame() string
}

type Animation struct {
	// Config part
	Framerate      float64
	Source         FrameSource
	Output         io.Writer
	ClearLastFrame bool

	duration  time.Duration
	init      sync.Once
	prevFrame string
	fstop     futuristicStop
}

func (animation *Animation) erase() {
	n := utf8.RuneCountInString(animation.prevFrame)
	back := strings.Repeat("\b", n)
	spaces := strings.Repeat(" ", n)
	fmt.Fprintf(animation.Output, "%s%s%s", back, spaces, back)
}

func (animation *Animation) Run() {
	animation.initAnimation()
	source := animation.Source
	defer animation.fstop.Init()()
	ticker := time.NewTicker(animation.duration)
	defer ticker.Stop()
cycle:
	for {
		select {
		case <-animation.fstop.Done():
			break cycle
		default:
			if !source.Next() {
				break cycle
			}
			frame := source.Frame()
			animation.erase()
			fmt.Fprintf(animation.Output, "%s", frame)
			animation.prevFrame = frame
			time.Sleep(animation.duration)
		}
	}
	if animation.ClearLastFrame {
		animation.erase()
	} else {
		fmt.Fprintf(animation.Output, "\n")
	}
}

func (animation *Animation) Stop() {
	animation.fstop.Stop()
}

func (animation *Animation) initAnimation() {
	animation.init.Do(func() {
		if animation.Output == nil {
			animation.Output = os.Stdout
		}
		if animation.Framerate == 0 {
			animation.Framerate = 5
		}
		animation.duration = time.Duration(float64(time.Second) / animation.Framerate)
	})
}
