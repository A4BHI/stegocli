package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"stegocli/config"

	"github.com/spf13/cobra"
)

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
var color = "\x1b[38;2;128;0;0m"
var reset = "\x1b[0m"
var encodeCmd = &cobra.Command{
	Use:   "encode -i image.png -f secretfile",
	Short: "Embed a secret file into a PNG image",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// var Image, File string
		var err error
		if cmd.Flags().NFlag() == 0 {
			cmd.Help()
			log.Fatal("No flags provided.")
		}
		enc := config.Encode{}
		enc.Image, err = cmd.Flags().GetString("image")
		if err != nil {
			fmt.Print(err)
			return
		}

		enc.SecretFile, err = cmd.Flags().GetString("file")
		if err != nil {
			fmt.Println(err)
			return
		}

		if enc.Image == "" || enc.SecretFile == "" {
			cmd.Help()
			log.Fatal("Not enough arguments.")

		}

		fmt.Println("Image Path : ", enc.Image, "\nFile Path : ", enc.SecretFile)

	},
}

var decodeCmd = &cobra.Command{
	Use:   "decode -i secretimage.png -o outputfilename",
	Short: "Extract a secret file from a PNG image",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		if cmd.Flags().NFlag() == 0 {
			cmd.Help()
			log.Fatal("Not enough flags.")
		}
		dec := config.Decode{}
		var err error

		dec.Image, err = cmd.Flags().GetString("image")
		if err != nil {
			log.Fatal(err)
		}

		dec.OutputFile, err = cmd.Flags().GetString("output")
		if err != nil {
			log.Fatal(err)
		}

		if dec.Image == "" || dec.OutputFile == "" {
			cmd.Help()
			log.Fatal("Not enough arguments.")
		}

	},
}

func init() {
	encodeCmd.Flags().StringP("image", "i", "", "Path to the image.")
	encodeCmd.Flags().StringP("file", "f", "", "Path to the secret file.")
	rootcmd.AddCommand(encodeCmd)
	decodeCmd.Flags().StringP("image", "i", "", "path to the secret image.")
	decodeCmd.Flags().StringP("output", "o", "", "Path to save the decoded file")
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
