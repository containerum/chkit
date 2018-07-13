package animation

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
	sliceFrames := make([]string, 0, len(frames))
	for _, frame := range frames {
		sliceFrames = append(sliceFrames, string(frame))
	}
	return FramesFromSlice(sliceFrames)
}
