package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
)

func Encrypt(compressedData []byte, password string) (ciphertext []byte, nonce []byte, salt []byte) {

	salt = make([]byte, 16)

	key, err := pbkdf2.Key(sha256.New, password, salt, 10000, 32)
	if err != nil {
		fmt.Println(err)

	}
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err)
	}

	nonce = make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}

	ciphertext = gcm.Seal(nil, nonce, compressedData, nil)

	return ciphertext, nonce, salt

}
