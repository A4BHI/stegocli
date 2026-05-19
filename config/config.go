package config

import (
	"time"

	"github.com/briandowns/spinner"
)

type Config struct {
	InputImage   string
	SecretFile   string
	OutputImage  string
	EncodedImage string
	DecodedFile  string
	Password     string
}

func StylenCallFunctions(function func() any, suffix string, finalmsg string) any {
	s := spinner.New(spinner.CharSets[14], 80*time.Millisecond)
	s.Suffix = suffix
	s.Color("cyan", "bold")
	s.FinalMSG = finalmsg
	s.Start()
	// data := function
	results := function()

	time.Sleep(1 * time.Second)

	s.Stop()

	return results
}
