package ui

import "fmt"

// PrintStatus prints a colored, emoji-prefixed status message.
func PrintStatus(icon, color, message string) {
	fmt.Printf("%s%s %s%s\n", color, icon, message, ColorReset)
}
