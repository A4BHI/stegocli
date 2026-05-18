package stego

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	_ "image/png"
	"path/filepath"

	"log"
	"os"
	"stegocli/compress"
	"stegocli/config"
	"stegocli/crypto"
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

	// s := spinner.New(spinner.CharSets[14], 80*time.Millisecond) // 100s was likely a typo for 100ms	s.Start()
	// s.Suffix = "  Compressing the secret file..."
	// s.Color("cyan", "bold")
	// s.FinalMSG = "\x1b[32m✔\x1b[0m Data compressed successfully.\n"
	// s.Start()
	// data := compress.Compress(cfg.SecretFile)
	// time.Sleep(1 * time.Second)

	// s.Stop()

	result := config.StylenCallFunctions(func() any {
		return compress.Compress(cfg.SecretFile)
	}, "Compressing the secret file...", "\x1b[32m✔\x1b[0m Data compressed successfully.\n")

	encryptionResults := config.StylenCallFunctions(func() any {
		ciphertext, nonce, salt := crypto.Encrypt(result.([]byte), cfg.Password)
		return struct {
			ciphertext []byte
			nonce      []byte
			salt       []byte
		}{ciphertext, nonce, salt}
	}, "Encrypting the secret file using AES-256", "\x1b[32m✔\x1b[0m Encryption Process Completed.")
	encryption := encryptionResults.(struct {
		ciphertext []byte
		nonce      []byte
		salt       []byte
	})
	index := 0
	length := len(encryption.ciphertext)
	// fmt.Println("Encoded length:", len(ciphertext))
	ext := filepath.Ext(cfg.SecretFile)
	extdata := []byte(ext)

	// magic := []byte("a4bhi")

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
	index = 0
	// bitsPrinted := 0
	// fmt.Print("Encode bits: ")

	// for bitsPrinted < 32 {

	// 	fmt.Print(pixels[index] & 1)
	// 	bitsPrinted++
	// 	index++
	// }
	// fmt.Println()
	err = png.Encode(OutFile, rgba)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Succesfully encoded the data inside the image.")

}
