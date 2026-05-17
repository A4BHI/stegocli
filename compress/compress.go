package compress

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"os"
)

func Compress(file string) []byte {
	data, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	var b bytes.Buffer

	gzip, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(gzip, data)
	err = gzip.Close()
	if err != nil {
		log.Fatal(err)
	}
	return b.Bytes()

}
