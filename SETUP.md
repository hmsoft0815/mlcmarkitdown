# Setup Guide: mlc-markitdown

MLC MarkItDown is a Go-based server that leverages Microsoft's `markitdown` Python library to convert various file formats into Markdown. This guide explains how to set up the required environment.

## Prerequisites

- **Go**: 1.25 or higher
- **Python**: 3.11 or higher
- **Artifact Server**: A running `mlcartifact` server (optional, but requested by default).

---

## Method A: Docker (Recommended)

Using Docker is the easiest way as it bundles all system dependencies (like `ffmpeg` and `libmagic`) and Python packages automatically.

1. **Build the image**:
   ```bash
   docker build -t mlc-markitdown .
   ```

2. **Run the container**:
   ```bash
   docker run -d \
     --name markitdown-server \
     -p 9591:9591 \
     --add-host=host.docker.internal:host-gateway \
     -e ARTIFACT_GRPC_ADDR=host.docker.internal:9590 \
     mlc-markitdown
   ```
   *Note: Use `host.docker.internal` to let the container talk back to services running on your host machine (like the artifact server).*

---

## Method B: Local Installation (Manual)

If you prefer to run it natively, you must install the Python dependencies and system libraries manually.

### 1. Install System Dependencies

**On Ubuntu/Debian:**
```bash
sudo apt-get update && sudo apt-get install -y \
    ffmpeg \
    libmagic1 \
    libxml2 \
    libxslt1-dev \
    gcc \
    python3-dev
```

**On macOS (Homebrew):**
```bash
brew install ffmpeg libmagic
```

### 2. Set up Python Environment

We recommend using a virtual environment:

```bash
# Create a virtual environment
python3 -m venv venv

# Activate it
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install the markitdown library
pip install markitdown
```

### 3. Build and Run the Go Server

1. **Build**:
   ```bash
   go build -o mlc-markitdown ./cmd/server/main.go
   ```

2. **Run**:
   Ensure your virtual environment is active so the Go server can find the `python3` command with `markitdown` installed.
   ```bash
   ./mlc-markitdown
   ```

---

## Configuration

You can configure the server using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | The port the server listens on | `9591` |
| `ARTIFACT_GRPC_ADDR` | Address of the mlcartifact server | `localhost:9590` |
| `PYTHON_CMD` | The command used to invoke Python | `python3` |
