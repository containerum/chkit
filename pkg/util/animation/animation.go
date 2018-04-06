package animation

import (
	"fmt"
	"io"
	"os"
	"runtime"
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
	isInitialised  bool
	done           chan struct{}
	stop           chan struct{}
}

func (animation *Animation) Run() {
	//ctx, stop := context.WithCancel(context.Background())
	animation.stop = make(chan struct{})
	animation.done = make(chan struct{})
	animation.initAnimation()
	source := animation.Source
	var prevFrame string
	if animation.ClearLastFrame {
		defer func(prevFrame *string) {
			fmt.Printf("\r%s\r", strings.Repeat(" ", UTF8stringWidth(*prevFrame)))
			select {
			case <-animation.done:
			default:
				close(animation.done)
			}
		}(&prevFrame)
	}
	if animation.WaitUntilDone != nil {
		go func() {
			animation.WaitUntilDone()
			select {
			case <-animation.stop:
			default:
				close(animation.stop)
			}
		}()
	}
	ticker := time.NewTicker(animation.duration)
	defer ticker.Stop()
	for {
		select {
		case <-animation.stop:
			return
		default:
			if !source.Next() {
				return
			}
			frame := source.Frame()
			fmt.Printf("\r%s\r%s", strings.Repeat(" ", UTF8stringWidth(prevFrame)), frame)
			prevFrame = frame
			time.Sleep(animation.duration)
		}
	}
}

func (animation *Animation) Stop() {
	for !animation.isInitialised {
		runtime.Gosched()
	}
	close(animation.stop)
	select {
	case <-animation.done:
	default:
		close(animation.done)
	}
}

func (animation *Animation) initAnimation() {
	defer func() { animation.isInitialised = true }()
	if animation.Output == nil {
		animation.Output = os.Stdout
	}
	if animation.Framerate == 0 {
		animation.Framerate = 5
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
