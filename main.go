package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"os/exec"
	"stegocli/config"
	"stegocli/stego"

	"github.com/spf13/cobra"
)

var color = "\x1b[38;2;128;0;0m"
var reset = "\x1b[0m"
var banner = color + `

  ██████ ▄▄▄█████▓▓█████   ▄████  ▒█████      ▄████▄   ██▓     ██▓
▒██    ▒ ▓  ██▒ ▓▒▓█   ▀  ██▒ ▀█▒▒██▒  ██▒   ▒██▀ ▀█  ▓██▒    ▓██▒
░ ▓██▄   ▒ ▓██░ ▒░▒███   ▒██░▄▄▄░▒██░  ██▒   ▒▓█    ▄ ▒██░    ▒██▒
  ▒   ██▒░ ▓██▓ ░ ▒▓█  ▄ ░▓█  ██▓▒██   ██░   ▒▓▓▄ ▄██▒▒██░    ░██░
▒██████▒▒  ▒██▒ ░ ░▒████▒░▒▓███▀▒░ ████▓▒░   ▒ ▓███▀ ░░██████▒░██░
▒ ▒▓▒ ▒ ░  ▒ ░░   ░░ ▒░ ░ ░▒   ▒ ░ ▒░▒░▒░    ░ ░▒ ▒  ░░ ▒░▓  ░░▓  
░ ░▒  ░ ░    ░     ░ ░  ░  ░   ░   ░ ▒ ▒░      ░  ▒   ░ ░ ▒  ░ ▒ ░
░  ░  ░    ░         ░   ░ ░   ░ ░ ░ ░ ▒     ░          ░ ░    ▒ ░
      ░              ░  ░      ░     ░ ░     ░ ░          ░  ░ ░  
                                             ░                    
` + reset

// var bold = "\x1b[44m"
var rootcmd = &cobra.Command{Use: "stego",
	Long: "A simple, user-friendly CLI that hides and extracts files in PNG images using LSB steganography.",
	Args: cobra.ArbitraryArgs,
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

var encodeCmd = &cobra.Command{
	Use:   "encode -i image.png -f secretfile -p password",
	Short: "Embed a secret file into a PNG image",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// var Image, File string
		var err error
		if cmd.Flags().NFlag() == 0 {
			cmd.Help()
			log.Fatal("No flags provided.")
		}
		enc := config.Config{}
		enc.InputImage, err = cmd.Flags().GetString("image")
		if err != nil {
			fmt.Print(err)
			return
		}

		enc.SecretFile, err = cmd.Flags().GetString("file")
		if err != nil {
			fmt.Println(err)
			return
		}

		enc.OutputImage, err = cmd.Flags().GetString("output")
		if err != nil {
			fmt.Println(err)
			return
		}

		enc.Password, err = cmd.Flags().GetString("password")
		if err != nil {
			fmt.Println(err)
			return
		}

		if enc.InputImage == "" || enc.SecretFile == "" || enc.Password == "" || enc.OutputImage == "" {
			cmd.Help()
			log.Fatal("Not enough arguments.")

		}

		if !Ispng(enc.InputImage) {
			log.Fatal("Tool supports .png images only.")
		}

		stego.Encode(&enc)

		// fmt.Println("Image Path : ", enc.Image, "\nFile Path : ", enc.SecretFile)

	},
}

var decodeCmd = &cobra.Command{
	Use:   "decode -i secretimage.png -o outputfilename -p password",
	Short: "Extract a secret file from a PNG image",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		if cmd.Flags().NFlag() == 0 {
			cmd.Help()
			log.Fatal("Not enough flags.")
		}
		dec := config.Config{}
		var err error

		dec.EncodedImage, err = cmd.Flags().GetString("image")
		if err != nil {
			log.Fatal(err)
		}

		dec.DecodedFile, err = cmd.Flags().GetString("output")
		if err != nil {
			log.Fatal(err)
		}

		dec.Password, err = cmd.Flags().GetString("password")
		if err != nil {
			log.Fatal(err)
		}

		if dec.EncodedImage == "" || dec.DecodedFile == "" || dec.Password == "" {
			cmd.Help()
			log.Fatal("Not enough arguments.")
		}

		stego.Decode(&dec)

	},
}

func init() {
	encodeCmd.Flags().StringP("image", "i", "", "Path to the image.")
	encodeCmd.Flags().StringP("file", "f", "", "Path to the secret file.")
	encodeCmd.Flags().StringP("output", "o", "", "Output image name with path.")
	encodeCmd.Flags().StringP("password", "p", "", "Password to encrypt the hidden data.")

	rootcmd.AddCommand(encodeCmd)
	decodeCmd.Flags().StringP("image", "i", "", "path to the secret image.")
	decodeCmd.Flags().StringP("output", "o", "", "Path to save the decoded file")
	decodeCmd.Flags().StringP("password", "p", "", "Password to decrypt the hidden data.")

	rootcmd.AddCommand(decodeCmd)
	rootcmd.CompletionOptions.DisableDefaultCmd = true

}

func main() {
	clear := exec.Command("clear")
	clear.Stdout = os.Stdout
	clear.Run()

	fmt.Println(banner)
	err := rootcmd.Execute()
	if err != nil {
		log.Fatal(err)
	}

}
