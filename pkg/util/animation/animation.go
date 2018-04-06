package animation

import (
	"context"
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

type sliceFrames struct {
	frames []string
	i      int
}

func (sl *sliceFrames) Next() bool {
	sl.i++
	if sl.i >= len(sl.frames) {
		sl.i = 0
	}
	return true
}

func (sl *sliceFrames) Frame() string {
	return sl.frames[sl.i]
}

func FramesFromSlice(frames []string) *sliceFrames {
	return &sliceFrames{
		i:      0,
		frames: frames,
	}
}

func FramesFromString(frames string) *sliceFrames {
	runes := []rune(frames)
	sliceFrames := make([]string, 0, len(frames))
	for _, frame := range runes {
		sliceFrames = append(sliceFrames, string(frame))
	}
	return FramesFromSlice(sliceFrames)
}

type Animation struct {
	Framerate      float64
	Source         FrameSource
	Output         io.Writer
	ClearLastFrame bool
	WaitUntilDone  func()
	duration       time.Duration
	init           sync.Once
	stop           func()
	done           chan struct{}
}

func (animation *Animation) Run() {
	animation.init.Do(func() { animation.initAnimation() })
	ctx, stop := context.WithCancel(context.Background())
	animation.done = make(chan struct{})
	defer func() { close(animation.done) }()
	animation.stop = stop
	defer stop()
	source := animation.Source
	var prevFrame string
	if animation.ClearLastFrame {
		defer func(prevFrame *string) {
			fmt.Fprintf(animation.Output, "\r%s\r", strings.Repeat(" ", UTF8stringWidth(*prevFrame)))
		}(&prevFrame)
	}
	if animation.WaitUntilDone != nil {
		go func() {
			animation.WaitUntilDone()
			stop()
		}()
	}
	ticker := time.NewTicker(animation.duration)
	defer ticker.Stop()
	for i := 0; true; i++ {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if !source.Next() {
				return
			}
			if i != 0 {
				fmt.Printf("\r%s\r", strings.Repeat(" ", UTF8stringWidth(prevFrame)))
			}
			frame := source.Frame()
			fmt.Printf("%s", frame)
			prevFrame = frame
		}
	}
}

func (animation *Animation) Stop() {
	if animation.stop != nil {
		animation.stop()
	}
	<-animation.done
}

func (animation *Animation) initAnimation() {
	if animation.Output == nil {
		animation.Output = os.Stdout
	}
	if animation.Framerate == 0 {
		animation.Framerate = 5
	}
	if animation.stop == nil {
		animation.stop = func() {}
	}
	if animation.done == nil {
		animation.done = make(chan struct{})
		close(animation.done)
	}
	animation.duration = time.Duration(float64(time.Second) / animation.Framerate)
}

func UTF8stringWidth(str string) int {
	size := 0
	for len(str) > 0 {
		_, n := utf8.DecodeRuneInString(str)
		str = str[n:]
		size += n
	}
	return size
}
