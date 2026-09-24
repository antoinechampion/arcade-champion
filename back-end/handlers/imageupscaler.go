package handlers

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"os/exec"
	"path/filepath"
)

type imageSize struct {
	width  int
	height int
}

var imageTargets = map[string]imageSize{
	"cover":  {width: 200, height: 267},
	"banner": {width: 1280, height: 720},
}

func upscaleImage(data []byte, kind, upscalerPath string) ([]byte, error) {
	if upscalerPath == "" {
		return data, nil
	}

	config, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode image dimensions: %w", err)
	}
	target := imageTargets[kind]
	if config.Width >= target.width && config.Height >= target.height {
		return data, nil
	}

	input, err := os.CreateTemp("", "arcade-champion-input-*.jpg")
	if err != nil {
		return nil, err
	}
	inputPath := input.Name()
	defer os.Remove(inputPath)
	outputPath := filepath.Join(filepath.Dir(inputPath), "arcade-champion-output-"+filepath.Base(inputPath))
	defer os.Remove(outputPath)

	if _, err := input.Write(data); err != nil {
		input.Close()
		return nil, err
	}
	if err := input.Close(); err != nil {
		return nil, err
	}

	cmd := exec.Command(upscalerPath, "-i", inputPath, "-o", outputPath, "-n", "realesrgan-x4plus")
	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("run image upscaler: %w: %s", err, output)
	}
	result, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, fmt.Errorf("read upscaled image: %w", err)
	}
	return result, nil
}
