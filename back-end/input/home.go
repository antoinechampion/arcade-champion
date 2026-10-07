package input

import (
	"encoding/binary"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	evKey     = 1
	btnMode   = 316 // gamepad "Home" button
	eventSize = 24  // struct input_event on 64-bit Linux
)

// WatchHomeHold calls onHold whenever a gamepad's Home button is held for `hold`.
// It reads /dev/input directly because the front-end gets no gamepad events while a game has focus.
func WatchHomeHold(hold time.Duration, onHold func()) {
	paths, _ := filepath.Glob("/dev/input/event*")
	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			log.Printf("[input] skipping %s: %v", path, err)
			continue
		}
		go watch(f, hold, onHold)
	}
}

func watch(r io.Reader, hold time.Duration, onHold func()) {
	var timer *time.Timer
	buf := make([]byte, eventSize)
	for {
		if _, err := io.ReadFull(r, buf); err != nil {
			return
		}
		typ := binary.LittleEndian.Uint16(buf[16:])
		code := binary.LittleEndian.Uint16(buf[18:])
		value := int32(binary.LittleEndian.Uint32(buf[20:]))
		if typ != evKey || code != btnMode {
			continue
		}
		switch value {
		case 1:
			timer = time.AfterFunc(hold, onHold)
		case 0:
			if timer != nil {
				timer.Stop()
			}
		}
	}
}
