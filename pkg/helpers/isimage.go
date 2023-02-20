package helpers

import (
	"mime/multipart"
	"net/http"
)

// validate if file is image
func IsImage(file multipart.File) bool {
	// get file header
	header := make([]byte, 512)
	_, err := file.Read(header)
	if err != nil {
		return false
	}
	// reset file pointer
	_, err = file.Seek(0, 0)
	if err != nil {
		return false
	}
	// validate if file is image
	return http.DetectContentType(header) == "image/jpg" || http.DetectContentType(header) == "image/png"
}

// is jpeg
func IsJpeg(file multipart.File) bool {
	// get file header
	header := make([]byte, 512)
	_, err := file.Read(header)
	if err != nil {
		return false
	}
	// reset file pointer
	_, err = file.Seek(0, 0)
	if err != nil {
		return false
	}
	// validate if file is image
	return http.DetectContentType(header) == "image/jpeg"
}
