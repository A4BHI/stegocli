package stego

import (
	"encoding/binary"
	"image"
	_ "image/png"
	"log"
	"os"

	"stegocli/compress"
	"stegocli/config"
	"stegocli/crypto"
)

type FileMetaData struct {
	Datalength int
	Extlength  int
	Extname    string
	CurrIndex  int
}

const (
	DATA_BIT_LENGTH       = 32
	FILE_EXTENSION_LENGTH = 8
	SALT_BITS             = 128
	NONCE_BITS            = 96
)

func Decode(cfg *config.Config) {
	inputimg, err := os.Open(cfg.EncodedImage + ".png")
	if err != nil {
		log.Fatal(err)
	}

	defer inputimg.Close()

	img, _, err := image.Decode(inputimg)
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
	filemetadata := FileMetaData{}

	datalen := readBytes(&filemetadata, pixels, DATA_BIT_LENGTH)
	lengthBytes := binary.BigEndian.Uint32(datalen)
	filemetadata.Datalength = int(lengthBytes) * 8

	extensionLength := readBytes(&filemetadata, pixels, FILE_EXTENSION_LENGTH)
	filemetadata.Extlength = int(extensionLength[0]) * 8

	filemetadata.Extname = string(readBytes(&filemetadata, pixels, filemetadata.Extlength))

	salt := readBytes(&filemetadata, pixels, SALT_BITS)
	nonce := readBytes(&filemetadata, pixels, NONCE_BITS)

	ciphertext := readBytes(&filemetadata, pixels, filemetadata.Datalength)

	plaintext := config.StylenCallFunctions(func() any {
		plaintext := crypto.Decrypt(ciphertext, salt, nonce, cfg.Password)
		return plaintext
	}, "\x1b[38;2;128;0;0mDecrypting embedded data.", "\x1b[32m✔\x1b[0m \x1b[38;2;0;255;0mDecryption completed.")

	DecompressedData := config.StylenCallFunctions(func() any {
		decompressedData := compress.Decompress(plaintext.([]byte))
		return decompressedData
	}, "\x1b[38;2;128;0;0mDecompressing extracted data.", "\x1b[32m✔\x1b[0m \x1b[38;2;0;255;0mDecompression completed.")

	config.StylenCallFunctions(func() any {
		err = os.WriteFile(cfg.DecodedFile+filemetadata.Extname, DecompressedData.([]byte), 0644)
		if err != nil {
			log.Fatal("save to device", err)
		}
		return nil
	}, "\x1b[38;2;128;0;0mWriting extracted data to disk...", "\x1b[32m✔\x1b[0m \x1b[38;2;0;255;0mFile written successfully\n\x1b[32m✔\x1b[0m \x1b[38;2;0;255;0mSaved as: "+cfg.DecodedFile+filemetadata.Extname)

}
func readBytes(fmd *FileMetaData, pixels []uint8, length int) []byte {
	bitsRead := 0
	var currbyte byte
	var byteslice []byte
	bitcount := 0
	for bitsRead < length {
		if fmd.CurrIndex >= len(pixels) {
			log.Fatal("payload exceeds image size")
		}
		bit := pixels[fmd.CurrIndex] & 1

		currbyte = (currbyte << 1) | bit
		bitcount++
		bitsRead++
		fmd.CurrIndex++
		if bitcount == 8 {
			byteslice = append(byteslice, currbyte)
			currbyte = 0
			bitcount = 0
		}
	}

	return byteslice
}
