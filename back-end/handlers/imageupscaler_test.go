package handlers

import (
	"bytes"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"
)

func TestUpscaleImageSkipsImagesAtTargetSize(t *testing.T) {
	data := jpegData(t, 200, 267)
	got, err := upscaleImage(data, "cover", "/does/not/run")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, data) {
		t.Fatal("expected image at target size to be unchanged")
	}
}

func TestUpscaleImageRunsConfiguredUpscaler(t *testing.T) {
	dir := t.TempDir()
	upscaler := filepath.Join(dir, "upscaler")
	script := "#!/bin/sh\ncp \"$2\" \"$4\"\n"
	if err := os.WriteFile(upscaler, []byte(script), 0755); err != nil {
		t.Fatal(err)
	}

	data := jpegData(t, 100, 100)
	got, err := upscaleImage(data, "cover", upscaler)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, data) {
		t.Fatal("expected test upscaler output")
	}
}

func jpegData(t *testing.T, width, height int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	var buffer bytes.Buffer
	if err := jpeg.Encode(&buffer, img, nil); err != nil {
		t.Fatal(err)
	}
	return buffer.Bytes()
}
