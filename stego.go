package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"os/exec"
	"stegocli/config"
	stego "stegocli/steganography"
	"strings"

	"github.com/spf13/cobra"
)

// var color = "\x1b[38;2;128;0;0m"
var color = "\x1b[38;2;255;85;85m"
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
	Long: "A Go-based CLI tool for embedding and extracting encrypted files or text payloads in PNG images using LSB steganography.",
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

var text, secretfile string

var encodeCmd = &cobra.Command{
	Use:   "encode -i image.png [-f file | -t text] -o outputimage.png -p password",
	Short: "Decode hidden files or text from a PNG image",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// var Image, File string
		var err error
		if cmd.Flags().NFlag() == 0 {
			cmd.Help()
			log.Fatal("No flags provided.")
		}
		enc := config.Config{}

		if enc.InputImage, err = cmd.Flags().GetString("image"); err != nil {
			log.Fatal(err)
		}

		if text, err = cmd.Flags().GetString("text"); err != nil {
			log.Fatal(err)

		}

		if secretfile, err = cmd.Flags().GetString("file"); err != nil {
			log.Fatal(err)
		}

		if text != "" {
			enc.SecretData = text
			enc.Flag = "text"
		}

		if secretfile != "" {
			enc.SecretData = secretfile
			enc.Flag = "file"
		}

		if text != "" && secretfile != "" || text == "" && secretfile == "" {
			cmd.Help()
			log.Fatal("Either use -t or use -f.")
		}

		if enc.OutputImage, err = cmd.Flags().GetString("output"); err != nil {
			log.Fatal(err)
		}

		if enc.Password, err = cmd.Flags().GetString("password"); err != nil {
			log.Fatal(err)
		}

		if enc.InputImage == "" || enc.Password == "" || enc.OutputImage == "" {
			cmd.Help()
			log.Fatal("Not enough arguments.")

		}

		if !Ispng(enc.InputImage) {
			log.Fatal("Tool supports .png images only.")
		}

		if !strings.HasSuffix(enc.OutputImage, ".png") {
			log.Fatal("Output file must end with .png")
		}

		stego.Encode(&enc)

	},
}

var decodeCmd = &cobra.Command{

	Use:   "decode -i secretimage.png  -p password",
	Short: "Extract hidden files or text from a PNG image",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		if cmd.Flags().NFlag() == 0 {
			cmd.Help()
			log.Fatal("Not enough flags.")
		}
		dec := config.Config{}
		var err error

		if dec.EncodedImage, err = cmd.Flags().GetString("image"); err != nil {
			log.Fatal(err)
		}

		if dec.Password, err = cmd.Flags().GetString("password"); err != nil {
			log.Fatal(err)
		}

		if dec.EncodedImage == "" || dec.Password == "" {
			cmd.Help()
			log.Fatal("Not enough arguments.")
		}
		if !Ispng(dec.EncodedImage) {
			log.Fatal("Tool supports .png images only.")
		}

		stego.Decode(&dec)

	},
}

func init() {
	encodeCmd.Flags().StringP("image", "i", "", "Path to the image.")
	encodeCmd.Flags().StringP("file", "f", "", "Path to the secret file.")
	encodeCmd.Flags().StringP("text", "t", "", "Directtly encode text into the image.")
	encodeCmd.Flags().StringP("output", "o", "", "Output image name with path.")
	encodeCmd.Flags().StringP("password", "p", "", "Password to encrypt the hidden data.")

	rootcmd.AddCommand(encodeCmd)
	decodeCmd.Flags().StringP("image", "i", "", "path to the secret image.")
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
