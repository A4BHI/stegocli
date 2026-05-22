package cmd

import (
	"log"
	"stegocli/config"
	stego "stegocli/steganography"
	"strings"

	"github.com/spf13/cobra"
)

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

		if !config.Ispng(enc.InputImage) {
			log.Fatal("Tool supports .png images only.")
		}

		if !strings.HasSuffix(enc.OutputImage, ".png") {
			log.Fatal("Output file must end with .png")
		}

		stego.Encode(&enc)

	},
}

func init() {
	encodeCmd.Flags().StringP("image", "i", "", "Path to the image.")
	encodeCmd.Flags().StringP("file", "f", "", "Path to the secret file.")
	encodeCmd.Flags().StringP("text", "t", "", "Directtly encode text into the image.")
	encodeCmd.Flags().StringP("output", "o", "", "Output image name with path.")
	encodeCmd.Flags().StringP("password", "p", "", "Password to encrypt the hidden data.")

	rootcmd.AddCommand(encodeCmd)
}
