#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Run go build
echo "Building passman-cli..."
go build -ldflags="-X 'github.com/NikitaAksenov/passman/internal/app.appConfiguration=prod'" -o passman.exe ./cmd/passman-cli

# Get the name of the binary (assumes current directory is the project root)
BINARY_NAME=$(basename "$PWD")

# Check if GOPATH is set
if [ -z "$GOPATH" ]; then
    echo "GOPATH is not set. Using default GOPATH..."
    GOPATH=$(go env GOPATH)
fi

# Move the binary to $GOPATH/bin
DEST="$GOPATH/bin/$BINARY_NAME"
echo "Moving binary to $DEST"
mv "$BINARY_NAME" "$DEST"

echo "Build and move completed successfully!"