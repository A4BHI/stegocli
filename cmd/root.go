package cmd

import "github.com/spf13/cobra"

var rootcmd = &cobra.Command{Use: "stego",
	Long: "A Go-based CLI tool for embedding and extracting encrypted files or text payloads in PNG images using LSB steganography.",
	Args: cobra.ArbitraryArgs,
}

func init() {
	rootcmd.CompletionOptions.DisableDefaultCmd = true
}

func Execute() {
	cobra.CheckErr(rootcmd.Execute())
}
