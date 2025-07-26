#!/usr/bin/env bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
GO_MOD_DIR="$SCRIPT_DIR/backend"
GO_TARGET_DIR="$SCRIPT_DIR/dist/go"
BUILD_BASE_NAME="chat_quest"
RUNTIME_UI_DIR="./browser" # NG dist relative to executable
RUNTIME_GIN_MODE="release"
BUILD_ARCHS=(
    "windows;amd64;${BUILD_BASE_NAME}_win_amd64.exe"
    "linux;amd64;${BUILD_BASE_NAME}_linux_amd64"
    "darwin;amd64;${BUILD_BASE_NAME}_mac_amd64"
)

cd "$GO_MOD_DIR"

for build_arch in "${BUILD_ARCHS[@]}"; do
    # Split the OS and ARCH from the array element
    IFS=';' read -r os arch filename <<< "$build_arch"
    target="$GO_TARGET_DIR/$filename"

    echo "Building ${os^} (${arch}) variant to '$target'..."
    GOOS=$os
    GOARCH=$arch
    go build \
      -ldflags="-X main.ChatQuestUIDir=$RUNTIME_UI_DIR -X main.GinMode=$RUNTIME_GIN_MODE" \
      -a -o "$target"
done
