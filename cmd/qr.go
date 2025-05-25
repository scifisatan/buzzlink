package cmd

import (
	"os"

	qrterminal "github.com/mdp/qrterminal/v3"
)

// ShowQR generates and prints a QR code in ANSI for the given link
func ShowQR(link string) error {
	config := qrterminal.Config{
		Level:      qrterminal.L,
		Writer:     os.Stdout,
		QuietZone:  1,    // Small border
		HalfBlocks: true, // Smaller QR code
	}
	qrterminal.GenerateWithConfig(link, config)
	return nil
}
