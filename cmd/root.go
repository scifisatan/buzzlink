package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)


var rootCmd = &cobra.Command{
	Use:   "buzzlink",
	Short: "buzzlink is a cli tool for sharing files and folders through buzzheavier's API",
	Long: `
	
    ██████╗ ██╗   ██╗███████╗███████╗██╗     ██╗███╗   ██╗██╗  ██╗
    ██╔══██╗██║   ██║╚══███╔╝╚══███╔╝██║     ██║████╗  ██║██║ ██╔╝
    ██████╔╝██║   ██║  ███╔╝   ███╔╝ ██║     ██║██╔██╗ ██║█████╔╝ 
    ██╔══██╗██║   ██║ ███╔╝   ███╔╝  ██║     ██║██║╚██╗██║██╔═██╗ 
    ██████╔╝╚██████╔╝███████╗███████╗███████╗██║██║ ╚████║██║  ██╗
    ╚═════╝  ╚═════╝ ╚══════╝╚══════╝╚══════╝╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝

Usage: 
  buzzlink [OPTIONS] <file|directory>

Options:
  -h        Show this help message
  -n NOTE   Add a note to the upload (optional)
  -u        Upgrade to the latest version from GitHub
  -p PASS   Password protect the upload before upload (optional)
  --noqr    Disable QR code display

Examples:
  buzzlink image.jpg                    # Upload a file
  buzzlink documents/                   # Upload a directory as zip
  buzzlink -n "Project files" src/     # Upload directory with note
  buzzlink -p "secret123" docs/        # Upload encrypted directory
  buzzlink image.jpg --noqr            # Upload without QR code
  
Requirements:
  • curl        - For file upload
  • xclip/wl-copy/pbcopy - For clipboard operations
  • 7z/zip      - For password protection

Report issues: github.com/scifisatan/buzzlink

	
	`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops, An error occured while executing Buzzlink, '%s' \n", err)
		os.Exit(1)
	}
}
