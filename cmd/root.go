package cmd

import (
	"fmt"
	"os"

	"buzzlink/cmd/archiver"
	"buzzlink/cmd/ui"

	"github.com/spf13/cobra"
	clipboard "github.com/tiagomelo/go-clipboard/clipboard"
)

var note string
var password string
var showQR bool

var rootCmd = &cobra.Command{
	Use:   "buzzlink",
	Short: "buzzlink is a cli tool for sharing files and folders through buzzheavier's API",
	Long:  ui.BuzzlinkBanner,
	Run: func(cmd *cobra.Command, args []string) {
		if note == "" && password == "" && !showQR && len(args) == 0 {
			fmt.Println(ui.BuzzlinkBanner)
			return
		}

		if len(args) == 0 {
			ui.PrintStatus(ui.IconError, ui.ColorRed, "No file or directory specified.")
			os.Exit(1)
		}

		filePath := args[0]
		link, err := handleFileUpload(filePath, note, password)
		if err != nil {
			ui.PrintStatus(ui.IconError, ui.ColorRed, fmt.Sprintf("Error: %v", err))
			os.Exit(1)
		}

		if showQR {
			ui.PrintStatus(ui.IconMobile, ui.ColorMagenta, "QR Code for mobile access:")
			ui.Show(link)
		}

		ui.PrintStatus(ui.IconLink, ui.ColorCyan, link)
		c := clipboard.New()
		if err := c.CopyText(link); err == nil {
			ui.PrintStatus(ui.IconClipboard, ui.ColorMagenta, "Copied to clipboard!")
		} else {
			ui.PrintStatus(ui.IconWarning, ui.ColorYellow, "Clipboard unavailable")
		}

		if password != "" {
			ui.PrintStatus(ui.IconKey, ui.ColorYellow, fmt.Sprintf("Password: %s (keep this safe!)", password))
		}
	},
}

func handleFileUpload(filePath, note, password string) (string, error) {
	zippedPath, err := archiver.ZipIfNeeded(filePath, password)
	if err != nil {
		return "", fmt.Errorf("zipping file: %w", err)
	}

	link, err := ui.RunUploadWithSpinner(zippedPath, note)
	if err != nil {
		if zippedPath != filePath {
			os.Remove(zippedPath)
		}
		return "", fmt.Errorf("uploading file: %w", err)
	}

	if zippedPath != filePath {
		os.Remove(zippedPath)
	}

	return link, nil
}

func init() {
	rootCmd.Flags().StringVarP(&note, "note", "n", "", "Add a note to the upload (optional)")
	rootCmd.Flags().StringVarP(&password, "password", "p", "", "Password protect the upload before upload (optional)")
	rootCmd.Flags().BoolVar(&showQR, "qr", false, "Enable QR code display")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		ui.PrintStatus(ui.IconError, ui.ColorRed, fmt.Sprintf("Oops, An error occurred while executing Buzzlink, '%s'", err))
		os.Exit(1)
	}
}
