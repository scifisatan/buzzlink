package ui

import (
	"os"

	qrterminal "github.com/mdp/qrterminal/v3"
)

// ShowQR generates and prints a QR code in ANSI for the given link.
// It's renamed to Show for brevity within the ui package context but was ShowQR.
func Show(link string) error { // Renamed from ShowQR, made public
	config := qrterminal.Config{
		Level:      qrterminal.L,
		Writer:     os.Stdout,
		QuietZone:  1,    // Small border
		HalfBlocks: true, // Smaller QR code
	}
	qrterminal.GenerateWithConfig(link, config)
	return nil
}
