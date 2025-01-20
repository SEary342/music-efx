#!/bin/bash

set -e

REPO="SEary342/music-efx"
APP_NAME="Music-EFX"
EXEC_NAME="music-efx"
ICON_URL="https://raw.githubusercontent.com/SEary342/music-efx/main/logo.png"

# Determine system paths
INSTALL_DIR="$HOME/.local/bin"
DESKTOP_FILE_DIR="$HOME/.local/share/applications"

# Function to fetch the latest release asset
download_latest_release() {
  echo "Fetching the latest release from $REPO..."
  API_URL="https://api.github.com/repos/$REPO/releases/latest"
  DOWNLOAD_URL=$(curl -s "$API_URL" | grep "browser_download_url.*$EXEC_NAME" | grep -v ".exe" | cut -d '"' -f 4)
  
  if [ -z "$DOWNLOAD_URL" ]; then
    echo "Error: Unable to fetch the latest release download URL." >&2
    exit 1
  fi

  echo "Downloading $EXEC_NAME from $DOWNLOAD_URL..."
  curl -L "$DOWNLOAD_URL" -o "$INSTALL_DIR/$EXEC_NAME"
  chmod +x "$INSTALL_DIR/$EXEC_NAME"
  echo "Executable installed at $INSTALL_DIR/$EXEC_NAME."
}

# Function to generate a .desktop file
generate_desktop_file() {
  echo "Generating .desktop file for $APP_NAME..."
  
  mkdir -p "$DESKTOP_FILE_DIR"
  
  cat > "$DESKTOP_FILE_DIR/$APP_NAME.desktop" <<EOF
[Desktop Entry]
Name=$APP_NAME
Exec=$INSTALL_DIR/$EXEC_NAME
Type=Application
Terminal=true
Categories=Utility;
EOF

  if [ -n "$ICON_URL" ]; then
    ICON_PATH="$DESKTOP_FILE_DIR/$APP_NAME.png"
    echo "Downloading icon from $ICON_URL..."
    curl "$ICON_URL" -o "$ICON_PATH"
    echo "Icon=$ICON_PATH" >> "$DESKTOP_FILE_DIR/$APP_NAME.desktop"
  fi

  chmod +x "$DESKTOP_FILE_DIR/$APP_NAME.desktop"
  echo "Desktop entry created at $DESKTOP_FILE_DIR/$APP_NAME.desktop."
}

# Ensure installation directory exists
mkdir -p "$INSTALL_DIR"

# Execute functions
download_latest_release
generate_desktop_file

echo "Installation complete. You can now find $APP_NAME in your application launcher."
