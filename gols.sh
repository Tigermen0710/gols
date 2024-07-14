#!/bin/bash

# Check if an argument is provided
if [ -n "$1" ]; then
    # Check if the directory exists
    if [ -d "$1" ]; then
        # Change to the specified directory
        cd "$1" || { echo "Failed to change directory to $1"; exit 1; }
    else
        echo "Directory $1 does not exist."
        exit 1
    fi
fi

# Execute the custom ls command
exec ~/gols/ls



