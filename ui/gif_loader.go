package ui

import (
	"os"
	"strings"
	"sync"
	"time"

	asciiconverter "github.com/mattparadis/asciiConverter"
)

var (
	gifFrames       []string
	gifFrameDelays  []time.Duration
	gifLoadFailed   bool
	loadedWidth     int
	gifLoadMu       sync.Mutex
	defaultGifWidth = 80
	gifSource       = "rocket.gif"
)

// loadGifFrames loads gears2.gif into colored ASCII frames at a target width.
// It caches the last successful load for the same width.
func loadGifFrames(width int) ([]string, []time.Duration, bool) {
	if width <= 0 {
		width = defaultGifWidth
	}

	if ok := ensureGifReady(); !ok {
		gifLoadFailed = true
		return nil, nil, false
	}

	gifLoadMu.Lock()
	defer gifLoadMu.Unlock()

	if width == loadedWidth && len(gifFrames) > 0 && !gifLoadFailed {
		return gifFrames, gifFrameDelays, true
	}

	frames, err := asciiconverter.GetAsciiGif(gifSource, width, 0)
	if err != nil {
		gifLoadFailed = true
		return nil, nil, false
	}

	gifFrames = gifFrames[:0]
	gifFrameDelays = gifFrameDelays[:0]
	for _, f := range frames {
		var sb strings.Builder
		for _, line := range f.Lines {
			sb.WriteString(line)
		}
		gifFrames = append(gifFrames, sb.String())
		gifFrameDelays = append(gifFrameDelays, f.Delay)
	}

	loadedWidth = width
	gifLoadFailed = false

	if gifLoadFailed || len(gifFrames) == 0 {
		return nil, nil, false
	}
	return gifFrames, gifFrameDelays, true
}

// chooseGifWidth picks a width that nearly fills the terminal, accounting for
// the converter writing two runes per pixel. Width is clamped to keep it usable.
func chooseGifWidth(termWidth int) int {
	if termWidth <= 0 {
		return defaultGifWidth
	}
	// Calculate width based on terminal size, leaving some margin
	w := (termWidth - 16) / 2
	if w < 40 {
		w = 40
	}
	if w > 100 {
		w = 100
	}
	return w
}

// ensureGifReady verifies the source GIF exists.
func ensureGifReady() bool {
	_, err := os.Stat(gifSource)
	return err == nil
}
