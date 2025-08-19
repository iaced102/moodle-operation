package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"

	"moodle/config"
)

// parse hex string in config.DB_ENC_KEY_HEX thành []byte key AES
func getSecretKey() []byte {
	keyBytes, err := hex.DecodeString(config.DB_ENC_KEY_HEX)
	if err != nil {
		log.Fatalf("Failed to decode DB_ENC_KEY_HEX: %v", err)
	}
	if len(keyBytes) != 32 {
		log.Fatalf("SecretKey must be 32 bytes (got %d)", len(keyBytes))
	}
	return keyBytes
}

func EncryptPassword(password string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(password))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(password))

	return hex.EncodeToString(ciphertext), nil
}

func DecryptPassword(encryptedHex string, key []byte) (string, error) {
	ciphertext, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func main() {
	mode := flag.String("mode", "", "encrypt or decrypt")
	text := flag.String("text", "", "the text to process")
	flag.Parse()

	if *mode == "" || *text == "" {
		log.Fatal("Usage: go run crypto.go -mode=<encrypt|decrypt> -text=<your_text>")
	}

	secretKey := getSecretKey()

	if *mode == "encrypt" {
		encrypted, err := EncryptPassword(*text, secretKey)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Encrypted:", encrypted)
	} else if *mode == "decrypt" {
		decrypted, err := DecryptPassword(*text, secretKey)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Decrypted:", decrypted)
	} else {
		log.Fatal("Invalid mode. Use encrypt or decrypt")
	}
}
