#!/bin/bash

# Configuration for Mac Studio Ollama
OLLAMA_URL="http://macstudio.fritz.box:11434/v1"
OLLAMA_MODEL="llava:7b"

echo "=== 1. Building MLC MarkItDown ==="
make build

echo ""
echo "=== 2. Running Vision & Model Test via Ollama ==="
echo "Target: $OLLAMA_URL ($OLLAMA_MODEL)"
echo ""

# We use the CLI flags to configure the server defaults for the test run
mcp-tester test \
  -c "bin/mlc-markitdown --llm-provider ollama --llm-url $OLLAMA_URL --llm-model $OLLAMA_MODEL" \
  --script tests/vision_test.mcp \
  -v

echo ""
echo "=== Test Finished ==="
