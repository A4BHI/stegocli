package stego

import (
	"encoding/binary"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
	"strings"

	"stegocli/compress"
	"stegocli/config"
	"stegocli/crypto"
)

type FileMetaData struct {
	Datalength     int
	FileNameLength int
	FileName       string
	CurrIndex      int
}

const (
	DATA_BITS        = 32
	FLAG_LENGTH      = 8
	FILE_NAME_LENGTH = 16
	SALT_BITS        = 128
	NONCE_BITS       = 96
	MODE_BITS        = 32
)

func Decode(cfg *config.Config) {
	if strings.Contains(cfg.EncodedImage, ".png") {

		splitted := strings.TrimRight(cfg.EncodedImage, ".png")
		cfg.EncodedImage = splitted
	}

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

	datalen := readBytes(&filemetadata, pixels, DATA_BITS)
	lengthBytes := binary.BigEndian.Uint32(datalen)
	filemetadata.Datalength = int(lengthBytes) * 8

	flaglen := readBytes(&filemetadata, pixels, FLAG_LENGTH)

	flag := string(readBytes(&filemetadata, pixels, int(flaglen[0])*8))
	var fileMode os.FileMode
	if flag == "file" {

		filenameLength := readBytes(&filemetadata, pixels, FILE_NAME_LENGTH)
		filemetadata.FileNameLength = int(
			binary.BigEndian.Uint16(filenameLength),
		)
		filemetadata.FileName = string(readBytes(&filemetadata, pixels, filemetadata.FileNameLength*8))

		fileMode = os.FileMode(binary.BigEndian.Uint32(readBytes(&filemetadata, pixels, MODE_BITS)))
	}

	salt := readBytes(&filemetadata, pixels, SALT_BITS)
	nonce := readBytes(&filemetadata, pixels, NONCE_BITS)

	ciphertext := readBytes(&filemetadata, pixels, filemetadata.Datalength)

	plaintext := config.StylenCallFunctions(func() any {
		plaintext := crypto.Decrypt(ciphertext, salt, nonce, cfg.Password)
		return plaintext
	}, "\x1b[38;2;255;85;85m Decrypting embedded data.", "\x1b[32m✔\x1b[0m \x1b[38;2;80;220;120mDecryption completed.")

	DecompressedData := config.StylenCallFunctions(func() any {
		decompressedData := compress.Decompress(plaintext.([]byte))
		return decompressedData
	}, "\x1b[38;2;255;85;85m Decompressing extracted data.", "\x1b[32m✔\x1b[0m \x1b[38;2;80;220;120mDecompression completed.")

	if flag == "file" {
		config.StylenCallFunctions(func() any {
			err = os.WriteFile(filemetadata.FileName, DecompressedData.([]byte), os.FileMode(fileMode))
			if err != nil {
				log.Fatal("save to device", err)
			}
			return nil
		}, "\x1b[38;2;255;85;85m Writing extracted data to disk...", "\x1b[32m✔\x1b[0m \x1b[38;2;120;190;255mFile written successfully\n\x1b[32m✔\x1b[0m \x1b[38;2;220;140;255mSaved as: "+filemetadata.FileName)

	} else {

		fmt.Println("\x1b[32m✔\x1b[0m \x1b[38;2;80;220;120mSecret extracted successfully\n" + "\x1b[38;2;120;190;255m\n[ Secret Data ]\n \x1b[38;2;255;215;0m➜\x1b[0m \x1b[38;2;220;140;255m" + string(DecompressedData.([]byte)))

	}

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
