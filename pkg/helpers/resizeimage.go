package helpers

import (
	"image"
	"image/jpeg"
	"os"

	"github.com/nfnt/resize"
)

// resize image
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



