package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/sha256"
	"fmt"
	"log"
)

func Decrypt(ciphertext []byte, salt []byte, nonce []byte, password string) []byte {

	key, err := pbkdf2.Key(sha256.New, password, salt, 10000, 32)
	if err != nil {
		log.Fatal(err)

	}
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err)
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File Decrypted")
	return plaintext
}
