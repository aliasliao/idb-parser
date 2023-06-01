#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "${0}")/.."
APP_ROOT=$(pwd)

function main() {
  MAIN="./main.go"

  if [[ "$#" -lt 1 ]]; then
    echo "Usage: $0 <w(indows)/d(arwin)>"
    exit 1
  elif [[ "$1" == "windows" ]] || [[ "$1" == "w" ]]; then
    export GOOS="windows"
    export GOARCH="amd64"
    EXT=".exe"
  elif [[ "$1" == "darwin" ]] || [[ "$1" == "d" ]] ; then
    export GOOS="darwin"
    export GOARCH="amd64"
    EXT=""
  elif [[ "$1" == "m1" ]] ; then
    export GOOS="darwin"
    export GOARCH="arm64"
    EXT=""
  fi

  BIN_PATH="$APP_ROOT/build/idb$EXT"
  go build -ldflags "-s -w" -o "$BIN_PATH" "$MAIN"
}

main "$@"
