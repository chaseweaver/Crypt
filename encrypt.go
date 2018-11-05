package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

// createHash(string)
// Creates an hashed string from a key
func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

// encryptAESFile([]byte, []byte)
// Encrypts a 64bit file
func encryptAESFile(code, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ct := make([]byte, aes.BlockSize+len(code))
	iv := ct[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ct[aes.BlockSize:], code)
	return ct, nil
}

// decdecryptAESFile([]byte, []byte)
// Decrpyts an encrypted 64bit file
func decryptAESFile(code, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Returns if the ciphertext is too short
	if len(code) < aes.BlockSize {
		return nil, nil
	}

	iv := code[:aes.BlockSize]
	code = code[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(code, code)
	return code, nil
}

// encryptMessage([]byte, string)
// Encrypts a 64bit string message
func encryptMessage(ext string, key []byte) (string, error) {
	pt := []byte(ext)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ct := make([]byte, aes.BlockSize+len(pt))
	iv := ct[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ct[aes.BlockSize:], pt)
	return base64.URLEncoding.EncodeToString(ct), nil
}

// decryptMessage([]byte, string)
// Decrpyts an encrypted 64bit message
func decryptMessage(ext string, key []byte) (string, error) {
	ct, _ := base64.URLEncoding.DecodeString(ext)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Returns if the ciphertext is too short
	if len(ct) < aes.BlockSize {
		return "", nil
	}

	iv := ct[:aes.BlockSize]
	ct = ct[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ct, ct)
	return fmt.Sprintf("%s", ct), nil
}

// returnMD5([]byte)
// Returns an MD5 hash of byte data
func returnMD5(data []byte) string {
	return fmt.Sprintf("%v", md5.Sum(data))
}
