package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

// Encrypt data using AES (CBC mode)
func Encrypt(plainText string, key []byte) (string, error) {
	// Generate a new random IV (Initialization Vector)
	iv := make([]byte, aes.BlockSize)
	_, err := io.ReadFull(rand.Reader, iv)
	if err != nil {
		return "", err
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Pad the plain text to a multiple of the AES block size
	paddedPlainText := pad([]byte(plainText), aes.BlockSize)

	// Encrypt the data using AES in CBC mode
	cipherText := make([]byte, len(paddedPlainText))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, paddedPlainText)

	// Prepend the IV to the cipher text to use during decryption
	encryptedData := append(iv, cipherText...)
	return hex.EncodeToString(encryptedData), nil
}

// Decrypt AES-encrypted data
func Decrypt(encryptedData string, key []byte) (string, error) {
	// Decode the hex string to get the encrypted data
	data, err := hex.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	// Extract the IV from the encrypted data (first block)
	iv := data[:aes.BlockSize]
	cipherText := data[aes.BlockSize:]

	// Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Decrypt the data using AES in CBC mode
	mode := cipher.NewCBCDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	mode.CryptBlocks(plainText, cipherText)

	// Unpad the decrypted data
	plainText = unpad(plainText)

	return string(plainText), nil
}

// Function to pad the plaintext to a multiple of the block size
func pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	paddingText := make([]byte, padding)
	for i := 0; i < padding; i++ {
		paddingText[i] = byte(padding)
	}
	return append(data, paddingText...)
}

// Function to remove padding from the decrypted plaintext
func unpad(data []byte) []byte {
	padding := int(data[len(data)-1])
	return data[:len(data)-padding]
}
