#!/bin/bash
#==============================================================================
#
#  ██████╗ ██╗   ██╗███████╗███████╗██╗     ██╗███╗   ██╗██╗  ██╗
#  ██╔══██╗██║   ██║╚══███╔╝╚══███╔╝██║     ██║████╗  ██║██║ ██╔╝
#  ██████╔╝██║   ██║  ███╔╝   ███╔╝ ██║     ██║██╔██╗ ██║█████╔╝ 
#  ██╔══██╗██║   ██║ ███╔╝   ███╔╝  ██║     ██║██║╚██╗██║██╔═██╗ 
#  ██████╔╝╚██████╔╝███████╗███████╗███████╗██║██║ ╚████║██║  ██╗
#  ╚═════╝  ╚═════╝ ╚══════╝╚══════╝╚══════╝╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝
#                                                                
#  📦 Simple file sharing utility
#==============================================================================

# Color definitions (using bright/bold versions for better contrast)
readonly COLOR_RED="\033[1;31m"
readonly COLOR_GREEN="\033[1;32m"
readonly COLOR_YELLOW="\033[1;33m"
readonly COLOR_BLUE="\033[1;36m"    # Using cyan instead of blue for better visibility
readonly COLOR_MAGENTA="\033[1;35m"
readonly COLOR_CYAN="\033[1;36m"
readonly COLOR_RESET="\033[0m"

# Icons
readonly ICON_SUCCESS="✅"
readonly ICON_ERROR="❌"
readonly ICON_WARNING="⚠️"
readonly ICON_INFO="ℹ️"
readonly ICON_LOCK="🔒"
readonly ICON_LINK="🔗"
readonly ICON_UPLOAD="📤"
readonly ICON_CLIPBOARD="📋"
readonly ICON_MOBILE="📱"
readonly ICON_KEY="🔑"
readonly ICON_PACKAGE="📦"

#==============================================================================
# Function Definitions
#==============================================================================

# Show help message
show_help() {
  cat << EOF
    ██████╗ ██╗   ██╗███████╗███████╗██╗     ██╗███╗   ██╗██╗  ██╗
    ██╔══██╗██║   ██║╚══███╔╝╚══███╔╝██║     ██║████╗  ██║██║ ██╔╝
    ██████╔╝██║   ██║  ███╔╝   ███╔╝ ██║     ██║██╔██╗ ██║█████╔╝ 
    ██╔══██╗██║   ██║ ███╔╝   ███╔╝  ██║     ██║██║╚██╗██║██╔═██╗ 
    ██████╔╝╚██████╔╝███████╗███████╗███████╗██║██║ ╚████║██║  ██╗
    ╚═════╝  ╚═════╝ ╚══════╝╚══════╝╚══════╝╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝

Usage: 
  $(basename "$0") [OPTIONS] <file|directory>

Options:
  -h        Show this help message
  -n NOTE   Add a note to the upload (optional)
  -p PASS   Password protect the upload before upload (optional)
  --noqr    Disable QR code display

Examples:
  $(basename "$0") image.jpg                    # Upload a file
  $(basename "$0") documents/                   # Upload a directory as zip
  $(basename "$0") -n "Project files" src/     # Upload directory with note
  $(basename "$0") -p "secret123" docs/        # Upload encrypted directory
  $(basename "$0") image.jpg --noqr            # Upload without QR code
  
Requirements:
  • curl        - For file upload
  • xclip/wl-copy/pbcopy - For clipboard operations
  • 7z/zip      - For password protection

Report issues: github.com/yourusername/buzzlink
EOF
  exit 0
}

# Print colored status messages
print_status() {
  local icon="$1"
  local color="$2"
  local message="$3"
  
  echo -e "${color}${icon} ${message}${COLOR_RESET}"
}

# Validate dependencies
check_dependencies() {
  local missing=()
  
  # Check for required tools
  command -v curl >/dev/null 2>&1 || missing+=("curl")
  
  # Check for at least one clipboard tool
  if ! command -v xclip >/dev/null 2>&1 && ! command -v wl-copy >/dev/null 2>&1 && ! command -v pbcopy >/dev/null 2>&1; then
    missing+=("xclip/wl-copy/pbcopy")
  fi
  
  # For password protection, check for either 7z or zip
  if ! command -v 7z >/dev/null 2>&1 && ! command -v zip >/dev/null 2>&1; then
    missing+=("7z/zip")
  fi
  
  if [ ${#missing[@]} -ne 0 ]; then
    print_status "$ICON_ERROR" "$COLOR_RED" "Missing required dependencies:"
    for dep in "${missing[@]}"; do
      echo -e "  • ${dep}"
    done
    exit 1
  fi
}

# Parse command line arguments
parse_args() {
  # Initialize variables
  FILE=""
  NOTE=""
  PASSWORD=""
  SHOW_QR=true

  local args=()
  
  # Separate standard arguments and --noqr flag
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --noqr)
        SHOW_QR=false
        ;;
      *)
        args+=("$1")
        ;;
    esac
    shift
  done
  
  # Set the positional arguments to our filtered args
  set -- "${args[@]}"
  
  # Check if any arguments were provided
  [[ $# -eq 0 ]] && { print_status "$ICON_ERROR" "$COLOR_RED" "No file specified."; show_help; }
  
  # Parse options using getopts
  while getopts ":hn:p:" opt; do
    case ${opt} in
      h)
        show_help
        ;;
      n)
        NOTE=$(echo -n "$OPTARG" | base64)
        ;;
      p)
        PASSWORD="$OPTARG"
        ;;
      \?)
        print_status "$ICON_ERROR" "$COLOR_RED" "Invalid option: -$OPTARG"
        show_help
        ;;
      :)
        print_status "$ICON_ERROR" "$COLOR_RED" "Option -$OPTARG requires an argument."
        show_help
        ;;
    esac
  done

  # Get the file from remaining arguments
  shift $((OPTIND - 1))
  
  # Check if a file was specified
  [[ $# -eq 0 ]] && { print_status "$ICON_ERROR" "$COLOR_RED" "No file specified."; show_help; }
  
  # Set and verify path exists
  FILE="$1"
  if [[ ! -e "$FILE" ]]; then
    print_status "$ICON_ERROR" "$COLOR_RED" "Path '$FILE' not found."
    exit 1
  fi
}

# Process file for upload (encryption if needed)
process_file() {
  local input_path="$1"
  local password="$2"
  local temp_dir
  local archive_name
  
  # Initialize variables
  UPLOAD_FILE="$input_path"
  
  # If it's a directory, create zip archive
  if [[ -d "$input_path" ]]; then
    # Create temporary directory
    TEMP_DIR=$(mktemp -d)
    archive_name="${TEMP_DIR}/$(basename "$input_path").zip"
    
    print_status "$ICON_PACKAGE" "$COLOR_BLUE" "Creating archive from directory..."
    
    # Try to create zip archive
    if command -v zip >/dev/null 2>&1; then
      (cd "$(dirname "$input_path")" && zip -r "$archive_name" "$(basename "$input_path")") >/dev/null 2>&1 || {
        print_status "$ICON_ERROR" "$COLOR_RED" "Failed to create archive from directory"
        rm -rf "$TEMP_DIR"
        exit 1
      }
      print_status "$ICON_SUCCESS" "$COLOR_GREEN" "Directory archived successfully"
      UPLOAD_FILE="$archive_name"
    else
      print_status "$ICON_ERROR" "$COLOR_RED" "zip is required for directory upload"
      rm -rf "$TEMP_DIR"
      exit 1
    fi
  fi
  
  # If password protection is requested
  if [[ -n "$password" ]]; then
    # Create temporary directory
    TEMP_DIR=$(mktemp -d)
    archive_name="${TEMP_DIR}/$(basename "$input_file").zip"
    
    print_status "$ICON_LOCK" "$COLOR_BLUE" "Creating encrypted archive..."
    
    # Try 7z first, fall back to zip
    if command -v 7z >/dev/null 2>&1; then
      7z a -p"$password" -mhe=on "$archive_name" "$input_file" >/dev/null 2>&1
      if [ $? -eq 0 ]; then
        print_status "$ICON_SUCCESS" "$COLOR_GREEN" "File encrypted successfully with 7z"
        UPLOAD_FILE="$archive_name"
      else
        print_status "$ICON_WARNING" "$COLOR_YELLOW" "7z encryption failed, trying zip..."
        if command -v zip >/dev/null 2>&1; then
          zip -j -P "$password" "$archive_name" "$input_file" >/dev/null 2>&1 || {
            print_status "$ICON_ERROR" "$COLOR_RED" "Failed to create encrypted archive with zip"
            rm -rf "$TEMP_DIR"
            exit 1
          }
          print_status "$ICON_SUCCESS" "$COLOR_GREEN" "File encrypted successfully with zip"
          UPLOAD_FILE="$archive_name"
        else
          print_status "$ICON_ERROR" "$COLOR_RED" "No encryption tools available (need 7z or zip)"
          rm -rf "$TEMP_DIR"
          exit 1
        fi
      fi
    else
      if command -v zip >/dev/null 2>&1; then
        zip -j -P "$password" "$archive_name" "$input_file" >/dev/null 2>&1 || {
          print_status "$ICON_ERROR" "$COLOR_RED" "Failed to create encrypted archive with zip"
          rm -rf "$TEMP_DIR"
          exit 1
        }
        print_status "$ICON_SUCCESS" "$COLOR_GREEN" "File encrypted successfully with zip"
        UPLOAD_FILE="$archive_name"
      else
        print_status "$ICON_ERROR" "$COLOR_RED" "No encryption tools available (need 7z or zip)"
        rm -rf "$TEMP_DIR"
        exit 1
      fi
    fi
  fi
  
  # Set filename and URL
  FILENAME=$(basename "$UPLOAD_FILE")
  UPLOAD_URL="https://w.buzzheavier.com/$FILENAME"
  [[ -n "$NOTE" ]] && UPLOAD_URL="${UPLOAD_URL}?note=$NOTE"
}

# Generate QR code for the download link using online API
generate_qr() {
  local link="$1"
  print_status "$ICON_MOBILE" "$COLOR_MAGENTA" "QR Code for mobile access:"
  echo
  
  # URL encode the link for safe transmission
  # Using standard URL encoding with bash/sed
  local encoded_link=$(echo "$link" | sed -e 's/%/%25/g' -e 's/ /%20/g' -e 's/!/%21/g' -e 's/"/%22/g' \
    -e 's/#/%23/g' -e 's/\$/%24/g' -e 's/\&/%26/g' -e "s/'/%27/g" -e 's/(/%28/g' \
    -e 's/)/%29/g' -e 's/\*/%2A/g' -e 's/+/%2B/g' -e 's/,/%2C/g' -e 's/-/%2D/g' \
    -e 's/\./%2E/g' -e 's/\//%2F/g' -e 's/:/%3A/g' -e 's/;/%3B/g' -e 's//%3C/g' \
    -e 's/=/%3D/g' -e 's/>/%3E/g' -e 's/?/%3F/g' -e 's/@/%40/g' -e 's/\[/%5B/g' \
    -e 's/\\/%5C/g' -e 's/\]/%5D/g' -e 's/\^/%5E/g' -e 's/_/%5F/g' -e 's/`/%60/g' \
    -e 's/{/%7B/g' -e 's/|/%7C/g' -e 's/}/%7D/g' -e 's/~/%7E/g')
  
  # Use QRCode.show API which directly returns ASCII art QR codes
  local qr_api_url="https://qrcode.show/$encoded_link?format=text"
  
  # Display the QR code in ASCII format
  if ! curl -s "$qr_api_url"; then
    print_status "$ICON_ERROR" "$COLOR_RED" "Failed to generate QR code from online service"
    return 1
  fi
  
  echo
}
# Copy link to clipboard
copy_to_clipboard() {
  local link="$1"
  if command -v xclip >/dev/null 2>&1; then
    echo -n "$link" | xclip -selection clipboard
    print_status "$ICON_CLIPBOARD" "$COLOR_CYAN" "Copied to clipboard (xclip)"
  elif command -v wl-copy >/dev/null 2>&1; then
    echo -n "$link" | wl-copy
    print_status "$ICON_CLIPBOARD" "$COLOR_CYAN" "Copied to clipboard (wl-copy)"
  elif command -v pbcopy >/dev/null 2>&1; then
    echo -n "$link" | pbcopy
    print_status "$ICON_CLIPBOARD" "$COLOR_CYAN" "Copied to clipboard (pbcopy)"
  else
    print_status "$ICON_WARNING" "$COLOR_YELLOW" "No clipboard tool found (xclip, wl-copy, or pbcopy). Install one manually."
    return 1
  fi
}

# Clean up temporary files
cleanup() {
  [[ -n "$TEMP_DIR" ]] && rm -rf "$TEMP_DIR"
}

#==============================================================================
# Main Script
#==============================================================================

# Set up cleanup on script exit
trap cleanup EXIT

# Check dependencies and parse arguments
check_dependencies
parse_args "$@"

# Process the file
process_file "$FILE" "$PASSWORD"

# Upload file and capture response with headers and verbose output
print_status "$ICON_UPLOAD" "$COLOR_BLUE" "Uploading $FILENAME..."
RESPONSE=$(curl -#o - -v -T "$UPLOAD_FILE" "$UPLOAD_URL" 2>&1)

# Extract JSON body (strip headers)
JSON_BODY=$(echo "$RESPONSE" | sed -n '/^{/,$p')

# Check if response contains JSON
if ! echo "$JSON_BODY" | grep -q '{'; then
  print_status "$ICON_ERROR" "$COLOR_RED" "Invalid response (not JSON)."
  exit 1
fi


# Use grep/cut to extract id from json
CODE=$(echo "$JSON_BODY" | grep -o '"code":[0-9]*' | cut -d':' -f2)
ID=$(echo "$JSON_BODY" | grep -o '"id":"[^"]*' | cut -d'"' -f4)


# Verify response code and ID
if [[ "$CODE" != "201" ]]; then
  print_status "$ICON_ERROR" "$COLOR_RED" "Upload failed with code $CODE."
  exit 1
fi

if [[ -z "$ID" ]]; then
  print_status "$ICON_ERROR" "$COLOR_RED" "Failed to extract ID from response."
  exit 1
fi

# Construct download link
LINK="https://buzzheavier.com/$ID"

# Output success
print_status "$ICON_SUCCESS" "$COLOR_GREEN" "Uploaded successfully!"

# Generate QR code if enabled
[[ "$SHOW_QR" = true ]] && generate_qr "$LINK"

# Show link and copy to clipboard
echo
print_status "$ICON_LINK" "$COLOR_CYAN" "$LINK"
echo

# Copy link to clipboard
copy_to_clipboard "$LINK"

# If file was password protected
[[ -n "$PASSWORD" ]] && print_status "$ICON_KEY" "$COLOR_YELLOW" "Password: $PASSWORD (keep this safe!)"

