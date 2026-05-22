package fightcade

import (
	"testing"
)

func TestComputeHash(t *testing.T) {
	hash := computeHash("tok123", "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", "6", "stcable", "abcdef0123456789abcdef0123456789")
	if len(hash) != 32 {
		t.Errorf("expected 32-char MD5 hex, got %d chars: %s", len(hash), hash)
	}
}

func TestReadCompressedUID_Invalid(t *testing.T) {
	got := readCompressedUID([]byte("not zlib"))
	if got != "" {
		t.Errorf("expected empty string for invalid data, got %q", got)
	}
}

func TestReadCompressedUID_Valid(t *testing.T) {
	uid := "12345678-1234-1234-1234-123456789abc"
	compressed, err := zlibCompress(uid)
	if err != nil {
		t.Fatal(err)
	}
	got := readCompressedUID(compressed)
	if got != uid {
		t.Errorf("got %q, want %q", got, uid)
	}
}
