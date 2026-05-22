package cmd

import (
	"log"
	"stegocli/config"
	stego "stegocli/steganography"

	"github.com/spf13/cobra"
)

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
		if !config.Ispng(dec.EncodedImage) {
			log.Fatal("Tool supports .png images only.")
		}

		stego.Decode(&dec)

	},
}

func init() {

	decodeCmd.Flags().StringP("image", "i", "", "path to the secret image.")
	decodeCmd.Flags().StringP("password", "p", "", "Password to decrypt the hidden data.")

	rootcmd.AddCommand(decodeCmd)
	rootcmd.CompletionOptions.DisableDefaultCmd = true

}
