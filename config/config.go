package config

import (
	"fmt"
	"image"
	"log"
	"os"
	"time"

	"github.com/briandowns/spinner"
)

type Config struct {
	InputImage   string
	SecretData   string
	OutputImage  string
	EncodedImage string
	Password     string
	Flag         string
}

func StylenCallFunctions(function func() any, suffix string, finalmsg string) any {

	s := spinner.New(spinner.CharSets[0], 80*time.Millisecond)
	s.Suffix = suffix
	s.Color("cyan", "bold")
	s.FinalMSG = finalmsg

	s.Start()

	results := function()

	time.Sleep(500 * time.Millisecond)

	s.Stop()
	fmt.Println()

	return results
}

func Ispng(imagepath string) bool {
	file, err := os.Open(imagepath)
	if err != nil {
		log.Fatal(err)
	}

	_, format, err := image.DecodeConfig(file)
	if err != nil {
		log.Fatal(err)
	}

	return format == "png"
}
