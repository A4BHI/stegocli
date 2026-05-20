package config

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

type Config struct {
	InputImage   string
	SecretFile   string
	OutputImage  string
	EncodedImage string
	// DecodedFile  string
	Password string
}

func StylenCallFunctions(function func() any, suffix string, finalmsg string) any {
	// spinnerset := []string{"S", "T", "E", "G", "O"}
	s := spinner.New(spinner.CharSets[0], 80*time.Millisecond)
	s.Suffix = suffix
	s.Color("cyan", "bold")
	s.FinalMSG = finalmsg

	s.Start()

	// data := function
	results := function()

	time.Sleep(1 * time.Second)

	s.Stop()
	fmt.Println()

	return results
}
