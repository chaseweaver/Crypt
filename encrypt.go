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
func encryptAESFile(code, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ct := make([]byte, aes.BlockSize+len(code))
	iv := ct[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ct[aes.BlockSize:], code)
	return ct
}

// decdecryptAESFile([]byte, []byte)
// Decrpyts an encrypted 64bit file
func decryptAESFile(code, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(code) < aes.BlockSize {
		panic("Ciphertext too short!")
	}
	iv := code[:aes.BlockSize]
	code = code[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(code, code)
	return code
}

// encryptExtension([]byte, string)
// Encrypts a 64bit string extension
func encryptExtension(ext string, key []byte) string {
	pt := []byte(ext)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ct := make([]byte, aes.BlockSize+len(pt))
	iv := ct[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ct[aes.BlockSize:], pt)
	return base64.URLEncoding.EncodeToString(ct)
}

// decryptExtension([]byte, string)
// Decrpyts an encrypted 64bit extension
func decryptExtension(ext string, key []byte) string {
	ct, _ := base64.URLEncoding.DecodeString(ext)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(ct) < aes.BlockSize {
		panic("Ciphertext too short!")
	}
	iv := ct[:aes.BlockSize]
	ct = ct[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ct, ct)
	return fmt.Sprintf("%s", ct)
}
