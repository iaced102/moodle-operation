package helpers

import (
	"encoding/base64"
)


// Encode encodes a struct to a base64 string
func Encode(v string) string {
	return base64.StdEncoding.EncodeToString([]byte(v))
}
