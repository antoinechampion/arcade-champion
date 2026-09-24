package handlers

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
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
	log.Printf("[image-upscaler] start kind=%q inputBytes=%d configured=%t path=%q", kind, len(data), upscalerPath != "", upscalerPath)
	if upscalerPath == "" {
		log.Printf("[image-upscaler] skip kind=%q reason=path-not-configured", kind)
		return data, nil
	}

	config, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		log.Printf("[image-upscaler] decode dimensions failed kind=%q inputBytes=%d error=%v", kind, len(data), err)
		return nil, fmt.Errorf("decode image dimensions: %w", err)
	}
	target := imageTargets[kind]
	log.Printf("[image-upscaler] dimensions kind=%q input=%dx%d target=%dx%d", kind, config.Width, config.Height, target.width, target.height)
	if config.Width >= target.width && config.Height >= target.height {
		log.Printf("[image-upscaler] skip kind=%q reason=already-at-target", kind)
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
	log.Printf("[image-upscaler] temp-files kind=%q input=%q output=%q", kind, inputPath, outputPath)

	if _, err := input.Write(data); err != nil {
		log.Printf("[image-upscaler] write input failed kind=%q path=%q error=%v", kind, inputPath, err)
		input.Close()
		return nil, err
	}
	if err := input.Close(); err != nil {
		log.Printf("[image-upscaler] close input failed kind=%q path=%q error=%v", kind, inputPath, err)
		return nil, err
	}

	log.Printf("[image-upscaler] execute kind=%q command=%q args=%q", kind, upscalerPath, []string{"-i", inputPath, "-o", outputPath, "-n", "realesrgan-x4plus"})
	cmd := exec.Command(upscalerPath, "-i", inputPath, "-o", outputPath, "-n", "realesrgan-x4plus")
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("[image-upscaler] execute failed kind=%q error=%v output=%q", kind, err, string(output))
		return nil, fmt.Errorf("run image upscaler: %w: %s", err, output)
	}
	log.Printf("[image-upscaler] execute succeeded kind=%q", kind)
	result, err := os.ReadFile(outputPath)
	if err != nil {
		log.Printf("[image-upscaler] read output failed kind=%q path=%q error=%v", kind, outputPath, err)
		return nil, fmt.Errorf("read upscaled image: %w", err)
	}
	resultConfig, resultFormat, resultErr := image.DecodeConfig(bytes.NewReader(result))
	if resultErr != nil {
		log.Printf("[image-upscaler] output validation failed kind=%q outputBytes=%d error=%v", kind, len(result), resultErr)
		return nil, fmt.Errorf("decode upscaled image: %w", resultErr)
	}
	log.Printf("[image-upscaler] complete kind=%q outputBytes=%d format=%q dimensions=%dx%d", kind, len(result), resultFormat, resultConfig.Width, resultConfig.Height)
	return result, nil
}
