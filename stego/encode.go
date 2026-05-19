package stego

import (
	"encoding/binary"
	"image"
	"image/png"
	_ "image/png"
	"path/filepath"
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
		return compress.Compress(cfg.SecretFile)
	}, "\x1b[38;2;128;0;0m Compressing the secret file...", "\x1b[32m✔\x1b[0m \x1b[38;2;0;255;0mData compressed successfully.")

	encryptionResults := config.StylenCallFunctions(func() any {
		ciphertext, nonce, salt := crypto.Encrypt(result.([]byte), cfg.Password)
		return struct {
			ciphertext []byte
			nonce      []byte
			salt       []byte
		}{ciphertext, nonce, salt}
	}, "\x1b[38;2;128;0;0mEncrypting the secret file using AES-256", "\x1b[32m✔\x1b[0m \x1b[38;2;0;255;0mEncryption Process Completed.")
	encryption := encryptionResults.(struct {
		ciphertext []byte
		nonce      []byte
		salt       []byte
	})

	s := spinner.New(spinner.CharSets[14], 80*time.Millisecond)
	s.Suffix = "\x1b[38;2;128;0;0mEmbedding encrypted payload into the carrier image..."
	s.Color("cyan", "bold")
	s.FinalMSG = "\x1b[32m✔\x1b[0m \x1b[38;2;0;255;0mSteganographic encoding complete. Payload secured.\n"
	s.Start()

	index := 0
	length := len(encryption.ciphertext)

	ext := filepath.Ext(cfg.SecretFile)
	extdata := []byte(ext)

	extbytes := byte(len(extdata))
	lengthBytes := make([]byte, 4)

	binary.BigEndian.PutUint32(lengthBytes, uint32(length))

	payload := append(lengthBytes, extbytes)
	payload = append(payload, extdata...)
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
	OutFile, err := os.Create(cfg.OutputImage + ".png")
	if err != nil {
		log.Fatal(err)
	}

	err = png.Encode(OutFile, rgba)
	if err != nil {
		log.Fatal(err)
	}

	s.Stop()

}
