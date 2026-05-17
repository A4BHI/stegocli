package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var rootcmd = &cobra.Command{Use: "stego",
	Long: "A basic steganography tool that can be used for encoding secrets data inside a png image.",
	Args: cobra.ArbitraryArgs,
	// Run:  rootFunc,
}

var encodeCmd = &cobra.Command{
	Use: "encode -i image.png -f secretfile",

	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var Image, File string
		var err error
		if cmd.Flags().NFlag() == 0 {
			cmd.Help()
			log.Fatal("No flags provided.")
		}

		Image, err = cmd.Flags().GetString("image")
		if err != nil {
			fmt.Print(err)
			return
		}

		File, err = cmd.Flags().GetString("file")
		if err != nil {
			fmt.Println(err)
			return
		}

		if Image == "" || File == "" {
			cmd.Help()
			log.Fatal("Not enough arguments.")

		}

		fmt.Println("Image Path : ", Image, "\nFile Path : ", File)

	},
}

var decodeCmd = &cobra.Command{
	Use: "decode -i secretimage.png -o outputfilename",
}

func init() {
	encodeCmd.Flags().StringP("image", "i", "", "Path to the image.")
	encodeCmd.Flags().StringP("file", "f", "", "Path to the secret file.")
	rootcmd.AddCommand(encodeCmd)
	rootcmd.AddCommand(decodeCmd)
}

func main() {

	err := rootcmd.Execute()
	if err != nil {
		log.Fatal(err)
	}

}
