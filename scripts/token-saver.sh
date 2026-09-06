#!/bin/bash

# Ensure user-installed binaries are in PATH
export PATH="$HOME/.local/bin:$PATH"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Check if graphify is installed
if ! command -v graphify &> /dev/null
then
    echo "graphify could not be found, installing..."
    PIP_FLAGS=""
    if pip install --help 2>&1 | grep -q "break-system-packages"; then
        PIP_FLAGS="--break-system-packages"
    fi
    pip install --user $PIP_FLAGS "graphifyy[gemini]"
fi

# Ensure the gemini backend dependency (openai) is installed
if ! python3 -c "import openai" &> /dev/null
then
    echo "openai package is missing, installing..."
    PIP_FLAGS=""
    if pip install --help 2>&1 | grep -q "break-system-packages"; then
        PIP_FLAGS="--break-system-packages"
    fi
    pip install --user $PIP_FLAGS openai
fi

echo "Running graphify extract on e2e-template..."
cd "$PROJECT_ROOT"
graphify extract . --backend gemini

