package fps

import (
	"log"
	"time"
)

// FrameRater ...
// Used to track speed of code
type FrameRater struct {
	startTime time.Time
	frames    int
	prepended string
}

// NewFrameRaterWithDescription ...
// When fps is found, preprend this string
func NewFrameRaterWithDescription(prepend string) FrameRater {
	return FrameRater{
		startTime: time.Now(),
		frames:    0,
		prepended: prepend,
	}
}

// NewFrameRater ...
// Generates new object
func NewFrameRater() FrameRater {
	return FrameRater{
		startTime: time.Now(),
		frames:    0,
		prepended: "Frame-Rater:",
	}
}

// Tick ...
// If more than one second has passed, the fps
// will be logged
func (f *FrameRater) Tick() {
	f.frames++
	now := time.Now()

	if diff := now.Sub(f.startTime); diff.Seconds() > 1.0 {
		log.Println(f.prepended, int(float64(f.frames)/diff.Seconds()), "fps")
		f.startTime = now
		f.frames = 0
	}
}
