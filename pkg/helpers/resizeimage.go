package helpers

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"strconv"

	"github.com/nfnt/resize"
)

// resize jpeg
func ResizeImage(imagePath string, width int, height int) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	m := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

	out, err := os.Create(imagePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)

	return imagePath, nil
}

// resize png
func ResizePng(imagePath string, width int, height int) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	m := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

	out, err := os.Create(imagePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// write new image to file
	png.Encode(out, m)

	return imagePath, nil
}


// resize ico
func ResizeIco(imagePath string, width int, height int) (string, error) {
	command := exec.Command("convert", imagePath, "-scale", strconv.Itoa(width)+"x"+strconv.Itoa(height), imagePath)
	err := command.Run()
	if err != nil {
		return "", err
	}

	return imagePath, nil
}
	

