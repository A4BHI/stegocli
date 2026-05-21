package compress

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"os"
	"stegocli/config"
)

func Compress(cfg *config.Config) []byte {

	var b bytes.Buffer
	var err error

	gzip, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	if err != nil {
		log.Fatal(err)
	}

	var file *os.File
	var data []byte

	if cfg.Flag == "file" {
		file, err = os.Open(cfg.SecretData)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		io.Copy(gzip, file)
		err = gzip.Close()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		data = []byte(cfg.SecretData)
		_, err = gzip.Write(data)
		if err != nil {
			log.Fatal(err)
		}

		err = gzip.Close()
		if err != nil {
			log.Fatal(err)
		}
	}

	return b.Bytes()

}
