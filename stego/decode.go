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

	// getFileExtension(&filemetadata, pixels)
	// salt, nonce := GetNonceandSalt(&filemetadata, pixels)

	salt := readBytes(&filemetadata, pixels, SALT_BITS)
	nonce := readBytes(&filemetadata, pixels, NONCE_BITS)
	// ciphertext := DecodeData(&filemetadata, pixels)
	ciphertext := readBytes(&filemetadata, pixels, filemetadata.Datalength)
	// fmt.Println("decoded ciphertext:", len(ciphertext))
	plaintext := crypto.Decrypt(ciphertext, salt, nonce, cfg.Password)
	// fmt.Println("decoded ciphertext:", len(ciphertext))
	DecodedData := compress.Decompress(plaintext)

	err = os.WriteFile(cfg.DecodedFile+filemetadata.Extname, DecodedData, 0644)
	if err != nil {
		log.Fatal("save to device", err)
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

// func getDatalenandExtLen(pixels []uint8) FileMetaData {
// 	index := 0
// 	bitsRead := 0

// 	var byteslice []byte
// 	var currbyte byte
// 	bitcount := 0
// 	var extlength byte
// 	for bitsRead < 40 {

// 		bit := pixels[index] & 1

// 		currbyte = (currbyte << 1) | bit
// 		bitcount++
// 		bitsRead++
// 		index++

// 		if bitcount == 8 && bitsRead <= 32 {
// 			byteslice = append(byteslice, currbyte)
// 			currbyte = 0
// 			bitcount = 0
// 		}

// 		if bitsRead >= 32 && bitcount == 8 {
// 			extlength = currbyte
// 			currbyte = 0
// 			bitcount = 0
// 		}
// 	}
// 	// fmt.Println(extlength)

// 	lengthBytes := binary.BigEndian.Uint32(byteslice)

// 	return FileMetaData{
// 		Datalength: int(lengthBytes) * 8,
// 		Extlength:  int(extlength) * 8,
// 		CurrIndex:  index,
// 	}
// }

// func getFileExtension(filemetadata *FileMetaData, pixels []uint8) {
// 	bitsRead := 0
// 	var currbyte byte
// 	bitcount := 0
// 	var sliceofext []byte
// 	for bitsRead < filemetadata.Extlength {
// 		bit := pixels[filemetadata.CurrIndex] & 1
// 		currbyte = (currbyte << 1) | bit
// 		filemetadata.CurrIndex++
// 		bitcount++
// 		bitsRead++

// 		if bitcount == 8 {
// 			sliceofext = append(sliceofext, currbyte)
// 			currbyte = 0
// 			bitcount = 0
// 		}
// 	}

// 	filemetadata.Extname = string(sliceofext)

// }
// func GetNonceandSalt(filemetadata *FileMetaData, pixels []uint8) ([]byte, []byte) {
// 	bitsRead := 0

// 	var saltslice []byte
// 	var nonceslice []byte

// 	var currbyte byte
// 	bitcount := 0

// 	for bitsRead < 224 { //salt * nonce means 16bytes * 12bytes = 224 bits we have to iterate till 224

// 		bit := pixels[filemetadata.CurrIndex] & 1

// 		currbyte = (currbyte << 1) | bit
// 		bitcount++
// 		bitsRead++
// 		filemetadata.CurrIndex++

// 		if bitcount == 8 && bitsRead <= 128 { //first 128 bits = salt
// 			saltslice = append(saltslice, currbyte)
// 			currbyte = 0
// 			bitcount = 0
// 		}

// 		if bitsRead > 128 && bitcount == 8 { //rest of the bits till 224 belongs to salt
// 			nonceslice = append(nonceslice, currbyte)
// 			currbyte = 0
// 			bitcount = 0
// 		}
// 	}

//		return saltslice, nonceslice
//	}
func DecodeData(filemetadata *FileMetaData, pixels []uint8) []byte {
	bitsRead := 0
	var sliceofdata []byte
	var currbyte byte
	bitcount := 0
	for bitsRead < filemetadata.Datalength {
		bit := pixels[filemetadata.CurrIndex] & 1
		currbyte = (currbyte << 1) | bit
		filemetadata.CurrIndex++
		bitcount++
		bitsRead++

		if bitcount == 8 {
			sliceofdata = append(sliceofdata, currbyte)
			currbyte = 0
			bitcount = 0
		}
	}
	// filename = filename + filemetadata.Extname
	// err := os.WriteFile(filename, sliceofdata, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return sliceofdata

}
