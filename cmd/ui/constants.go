package ui

// Color constants (ANSI escape codes)
const (
	ColorRed     = "\033[1;31m"
	ColorGreen   = "\033[1;32m"
	ColorYellow  = "\033[1;33m"
	ColorBlue    = "\033[1;36m" // Cyan for better visibility
	ColorMagenta = "\033[1;35m"
	ColorCyan    = "\033[1;36m"
	ColorReset   = "\033[0m"
)

// Emoji/icon constants
const (
	IconSuccess   = "✅"
	IconError     = "❌"
	IconWarning   = "⚠️"
	IconInfo      = "ℹ️"
	IconLock      = "🔒"
	IconLink      = "🔗"
	IconUpload    = "📤"
	IconClipboard = "📋"
	IconMobile    = "📱"
	IconKey       = "🔑"
	IconPackage   = "📦"
)

var BuzzlinkBanner = `
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
  --qr      Enable QR code display

For more information, visit: https://buzzlink.netify.app
Report issues: github.com/scifisatan/buzzlink
`
