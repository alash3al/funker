package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"io"

	"golang.org/x/crypto/bcrypt"
)

func jsCrypto() map[string]interface{} {
	return map[string]interface{}{
		"md5": func(s string) string {
			hash := md5.Sum([]byte(s))
			return hex.EncodeToString(hash[:])
		},
		"sha1": func(s string) string {
			hash := sha1.Sum([]byte(s))
			return hex.EncodeToString(hash[:])
		},
		"sha256": func(s string) string {
			hash := sha256.Sum256([]byte(s))
			return hex.EncodeToString(hash[:])
		},
		"sha512": func(s string) string {
			hash := sha512.Sum512([]byte(s))
			return hex.EncodeToString(hash[:])
		},
		"bcrypt": map[string]interface{}{
			"hash": func(s string) string {
				d, e := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
				if e != nil {
					panic(e.Error())
				}
				return hex.EncodeToString(d)
			},
			"check": func(hash, password string) bool {
				return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
			},
		},
		"aes": map[string]interface{}{
			"encrypt": func(plaintext, key string) string {
				c, err := aes.NewCipher([]byte(key))
				if err != nil {
					panic(err.Error())
				}

				gcm, err := cipher.NewGCM(c)
				if err != nil {
					panic(err.Error())
				}

				nonce := make([]byte, gcm.NonceSize())
				if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
					panic(err.Error())
				}

				return hex.EncodeToString(gcm.Seal(nonce, nonce, []byte(plaintext), nil))
			},
			"decrypt": func(ciphertext, key string) string {
				ciphertextB, err := hex.DecodeString(ciphertext)
				if err != nil {
					panic(err.Error())
				}
				keyB := []byte(key)
				c, err := aes.NewCipher(keyB)
				if err != nil {
					panic(err.Error())
				}

				gcm, err := cipher.NewGCM(c)
				if err != nil {
					panic(err.Error())
				}

				nonceSize := gcm.NonceSize()
				if len(ciphertextB) < nonceSize {
					panic("ciphertext too short")
				}

				nonce, ciphertextB := ciphertextB[:nonceSize], ciphertextB[nonceSize:]
				b, e := gcm.Open(nil, nonce, ciphertextB, nil)
				if e != nil {
					panic(e.Error())
				}
				return string(b)
			},
		},
	}
}
