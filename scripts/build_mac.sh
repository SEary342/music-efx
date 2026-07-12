#!/bin/bash
set -e

APP_NAME="MusicEFX"
DIST_DIR="dist"
APP_DIR="$DIST_DIR/$APP_NAME.app"
ICONSET_DIR="$DIST_DIR/$APP_NAME.iconset"
SOURCE_EXE="$DIST_DIR/music-efx"
LOGO_IMG="logoBig.png"
BUNDLE_ID="com.musicefx.player"

echo "🔨 Building macOS application bundle..."

# 1. Compile AppleScript App
osacompile -o "$APP_DIR" \
  -e 'set execPath to (POSIX path of (path to me)) & "Contents/Resources/music-efx"' \
  -e 'tell application "Terminal"' \
  -e 'activate' \
  -e 'do script quoted form of execPath' \
  -e 'end tell'

# 2. Copy the PyInstaller executable inside the bundle
mkdir -p "$APP_DIR/Contents/Resources"
cp "$SOURCE_EXE" "$APP_DIR/Contents/Resources/"

# 3. Remove the Asset Catalog so macOS uses our .icns file instead of the
#    default AppleScript applet icon baked into Assets.car
rm -f "$APP_DIR/Contents/Resources/Assets.car"

# 4. Build Icons
echo "🎨 Generating Retina-ready app icons..."
mkdir -p "$ICONSET_DIR"

sizes=(16 32 128 256 512)
for size in "${sizes[@]}"; do
    sips -z "$size" "$size" "$LOGO_IMG" --out "$ICONSET_DIR/icon_${size}x${size}.png" > /dev/null 2>&1
    double_size=$((size * 2))
    sips -z "$double_size" "$double_size" "$LOGO_IMG" --out "$ICONSET_DIR/icon_${size}x${size}@2x.png" > /dev/null 2>&1
done

iconutil -c icns "$ICONSET_DIR" -o "$APP_DIR/Contents/Resources/applet.icns"
rm -rf "$ICONSET_DIR"

# 5. Patch Info.plist: set bundle ID and remove CFBundleIconName so macOS
#    doesn't look for the icon in the (now-deleted) Asset Catalog
plutil -replace CFBundleIdentifier -string "$BUNDLE_ID" "$APP_DIR/Contents/Info.plist"
plutil -remove CFBundleIconName "$APP_DIR/Contents/Info.plist" 2>/dev/null || true

# 6. Fix the broken signature by applying a fresh ad-hoc signature
echo "🔐 Ad-hoc signing the application bundle..."
codesign --force --deep --sign - "$APP_DIR"

# 7. Handle Installation
if [ "$1" == "--install" ]; then
    echo "🚚 Installing to /Applications..."
    rm -rf "/Applications/$APP_NAME.app"
    cp -R "$APP_DIR" "/Applications/"

    echo "🔍 Force-refreshing icon and Spotlight caches..."
    # Reset Launch Services database so Finder picks up the new icon
    /System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister \
        -f "/Applications/$APP_NAME.app"
    # Touch the app so Finder knows it changed
    touch "/Applications/$APP_NAME.app"
    # Re-index for Spotlight
    mdimport "/Applications/$APP_NAME.app"
    # Restart Dock to flush icon caches (also restarts Finder icon rendering)
    killall Dock 2>/dev/null || true

    echo "🎉 Successfully installed $APP_NAME.app to your Applications folder!"
else
    echo "✨ Successfully built $APP_NAME.app in $DIST_DIR/"
fi