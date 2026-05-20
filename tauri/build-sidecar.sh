#!/bin/sh
set -e

TARGET_TRIPLE=$(rustc -vV | grep '^host:' | cut -d' ' -f2)
OUT="tauri/binaries/arcade-champion-backend-${TARGET_TRIPLE}"

echo "Building Go backend → ${OUT}"
cd back-end
go build -o "../${OUT}" .
