#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Run go build
echo "Building passman-desktop..."
go build -ldflags="-X 'github.com/NikitaAksenov/passman/internal/app.appConfiguration=prod'" -o passman.exe ./cmd/passman-desktop

echo "Build passman-desktop completed successfully!"