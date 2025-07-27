#!/usr/bin/env bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

DO_BUILD_UI=true
DO_BUILD_GO_BACKEND=true
DO_CREATE_PACKAGES=true

DIST_DIR="$SCRIPT_DIR/dist"
NG_BROWSER_DIST_DIR="$DIST_DIR/chat-quest/browser"
GO_MOD_DIR="$SCRIPT_DIR/backend"
GO_TARGET_DIR="$DIST_DIR/go"
GO_BUILD_BASE_NAME="chat_quest"
GO_RUNTIME_UI_DIR="browser" # NG dist relative to executable
GO_RUNTIME_GIN_MODE="release"
GO_BUILD_ARCHS=(
    "windows;amd64;${GO_BUILD_BASE_NAME}_win_amd64.exe"
    "linux;amd64;${GO_BUILD_BASE_NAME}_linux_amd64"
    "darwin;amd64;${GO_BUILD_BASE_NAME}_mac_amd64"
)

if [[ "$DO_BUILD_UI" == "true" ]]; then
  echo "Building UI..."
  npm run build
fi

if [[ "$DO_BUILD_GO_BACKEND" == "true" ]]; then
  echo ""
  echo "Building Go backend..."
  cd "$GO_MOD_DIR"
  for build_arch in "${GO_BUILD_ARCHS[@]}"; do
      # Split the OS and ARCH from the array element
      IFS=';' read -r os arch filename <<< "$build_arch"
      target="$GO_TARGET_DIR/$filename"

      echo "Building ${os^} (${arch}) variant..."
      GOOS=$os
      GOARCH=$arch
      go build \
        -ldflags="-X main.ChatQuestUIDir=$GO_RUNTIME_UI_DIR -X main.GinMode=$GO_RUNTIME_GIN_MODE" \
        -a -o "$target"
  done
fi

if [[ "$DO_CREATE_PACKAGES" == "true" ]]; then
  echo ""
  echo "Creating distributable packages..."
  rm -rf "$DIST_DIR/*.zip"

  for build_arch in "${GO_BUILD_ARCHS[@]}"; do
      IFS=';' read -r os arch filename <<< "$build_arch"

      # Create a temporary directory for zipping
      temp_dir="$DIST_DIR/temp_${filename%.*}"
      mkdir -p "$temp_dir"
      cp "$GO_TARGET_DIR/$filename" "$temp_dir/"

      # Copy the Angular distribution files
      mkdir -p "$temp_dir/browser"
      cp -r "$NG_BROWSER_DIST_DIR" "$temp_dir/"

      # Create the ZIP file
      zip_file="$DIST_DIR/${filename%.*}.zip"
      (cd "$temp_dir" && zip -qr "$zip_file" ./*)

      # Clean up the temporary directory
      rm -rf "$temp_dir"

      echo "Created distribution package: $zip_file"
  done
fi

echo "Build completed successfully!"
