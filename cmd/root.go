package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	clipboard "github.com/tiagomelo/go-clipboard/clipboard"
)

var note string
var password string
var showQR bool

// printStatus prints a colored, emoji-prefixed status message
func printStatus(icon, color, message string) {
	fmt.Printf("%s%s %s%s\n", color, icon, message, ColorReset)
}

var rootCmd = &cobra.Command{
	Use:   "buzzlink",
	Short: "buzzlink is a cli tool for sharing files and folders through buzzheavier's API",
	Long:  BuzzlinkBanner,
	Run: func(cmd *cobra.Command, args []string) {
		if note == "" && password == "" && !showQR && len(args) == 0 {
			fmt.Println(BuzzlinkBanner)
			return
		}

		if len(args) == 0 {
			printStatus(IconError, ColorRed, "No file or directory specified.")
			os.Exit(1)
		}

		filePath := args[0]

		// Step 1: Zip if needed
		zippedPath, err := ZipIfNeeded(filePath, password)
		if err != nil {
			printStatus(IconError, ColorRed, fmt.Sprintf("Error zipping file: %v", err))
			os.Exit(1)
		}

		// Step 2: Upload with spinner
		link, err := runUploadWithSpinner(zippedPath, note)
		if err != nil {
			printStatus(IconError, ColorRed, fmt.Sprintf("Error uploading file: %v", err))
			if zippedPath != filePath {
				os.Remove(zippedPath)
			}
			os.Exit(1)
		}

		// Step 3: Show QR if requested
		if showQR {
			printStatus(IconMobile, ColorMagenta, "QR Code for mobile access:")
			ShowQR(link)
		}

		// Step 4: Print download link and copy to clipboard (after QR)
		// printStatus(IconSuccess, ColorGreen, "Uploaded successfully!")
		printStatus(IconLink, ColorCyan, link)
		c := clipboard.New()
		if err := c.CopyText(link); err == nil {
			printStatus(IconClipboard, ColorCyan, "Copied to clipboard!")
		} else {
			printStatus(IconWarning, ColorYellow, "Clipboard unavailable")
		}

		// Step 5: Show password if set
		if password != "" {
			printStatus(IconKey, ColorYellow, fmt.Sprintf("Password: %s (keep this safe!)", password))
		}

		// Step 6: Cleanup zip if generated
		if zippedPath != filePath {
			os.Remove(zippedPath)
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&note, "note", "n", "", "Add a note to the upload (optional)")
	rootCmd.Flags().StringVarP(&password, "password", "p", "", "Password protect the upload before upload (optional)")
	rootCmd.Flags().BoolVar(&showQR, "qr", false, "Enable QR code display")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops, An error occured while executing Buzzlink, '%s' \n", err)
		os.Exit(1)
	}
}
