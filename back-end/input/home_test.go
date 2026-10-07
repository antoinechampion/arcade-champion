package input

import (
	"encoding/binary"
	"io"
	"testing"
	"time"
)

func event(typ, code uint16, value int32) []byte {
	b := make([]byte, eventSize)
	binary.LittleEndian.PutUint16(b[16:], typ)
	binary.LittleEndian.PutUint16(b[18:], code)
	binary.LittleEndian.PutUint32(b[20:], uint32(value))
	return b
}

// run feeds events to watch, then waits past the hold duration and reports whether onHold fired.
func run(t *testing.T, events ...[]byte) bool {
	t.Helper()
	pr, pw := io.Pipe()
	fired := make(chan struct{}, 1)
	go watch(pr, 50*time.Millisecond, func() { fired <- struct{}{} })
	for _, e := range events {
		pw.Write(e)
	}
	time.Sleep(150 * time.Millisecond)
	pw.Close()
	select {
	case <-fired:
		return true
	default:
		return false
	}
}

func TestHomeHeldTriggers(t *testing.T) {
	if !run(t, event(evKey, btnMode, 1)) {
		t.Fatal("expected onHold after holding Home")
	}
}

func TestHomeTapDoesNotTrigger(t *testing.T) {
	if run(t, event(evKey, btnMode, 1), event(evKey, btnMode, 0)) {
		t.Fatal("release before hold duration must not trigger")
	}
}

func TestOtherButtonsIgnored(t *testing.T) {
	if run(t, event(evKey, 304, 1), event(evKey, 304, 1)) {
		t.Fatal("non-Home button must not trigger")
	}
}
