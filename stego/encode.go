package stego

import (
	"encoding/binary"
	"image"
	"image/png"
	_ "image/png"
	"path/filepath"
	"strings"
	"time"

	"log"
	"os"
	"stegocli/compress"
	"stegocli/config"
	"stegocli/crypto"

	"github.com/briandowns/spinner"
)

func Encode(cfg *config.Config) {
	file, err := os.Open(cfg.InputImage)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	bounds := img.Bounds()
	rgba := image.NewNRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}

	pixels := rgba.Pix

	result := config.StylenCallFunctions(func() any {
		return compress.Compress(cfg.SecretData)
	}, "\x1b[38;2;255;85;85m Compressing the secret file...", "\x1b[32m✔\x1b[0m \x1b[38;2;0;255;0mData compressed successfully.")

	encryptionResults := config.StylenCallFunctions(func() any {
		ciphertext, nonce, salt := crypto.Encrypt(result.([]byte), cfg.Password)
		return struct {
			ciphertext []byte
			nonce      []byte
			salt       []byte
		}{ciphertext, nonce, salt}
	}, "\x1b[38;2;255;85;85m Encrypting the secret file using AES-256", "\x1b[32m✔\x1b[0m \x1b[38;2;0;255;0mEncryption Process Completed.")
	encryption := encryptionResults.(struct {
		ciphertext []byte
		nonce      []byte
		salt       []byte
	})

	s := spinner.New(spinner.CharSets[0], 80*time.Millisecond)
	s.Suffix = "\x1b[38;2;255;85;85m Embedding encrypted payload into the carrier image..."
	s.Color("cyan", "bold")
	s.FinalMSG = "\x1b[32m✔\x1b[0m \x1b[38;2;0;255;0mSteganographic encoding complete. Payload secured.\n"
	s.Start()

	index := 0
	length := len(encryption.ciphertext)
	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, uint32(length))

	flagvalue := []byte(cfg.Flag)
	flaglength := byte(len(flagvalue))

	payload := append(lengthBytes, flaglength)
	payload = append(payload, flagvalue...)

	if cfg.Flag == "file" {
		filename := filepath.Base(cfg.SecretData)
		nameData := []byte(filename)
		nameLenBytes := make([]byte, 2)
		binary.BigEndian.PutUint16(nameLenBytes, uint16(len(nameData)))

		payload = append(payload, nameLenBytes...)
		payload = append(payload, nameData...)
	}

	payload = append(payload, encryption.salt...)
	payload = append(payload, encryption.nonce...)
	payload = append(payload, encryption.ciphertext...)
	totalbits := len(payload) * 8

	if totalbits > len(pixels) {
		log.Fatal("Not enough space in image.")
	}
	for i := 0; i < len(payload); i++ {
		for j := 7; j >= 0; j-- {

			bit := (payload[i] >> j) & 1

			pixels[index] = (pixels[index] & 254) | bit
			index++
		}

	}
	var OutFile *os.File
	if strings.Contains(cfg.OutputImage, ".png") {
		OutFile, err = os.Create(cfg.OutputImage)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		OutFile, err = os.Create(cfg.OutputImage + ".png")
		if err != nil {
			log.Fatal(err)
		}
	}

	err = png.Encode(OutFile, rgba)
	if err != nil {
		log.Fatal(err)
	}

	s.Stop()

}
