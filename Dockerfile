# Stage 1: Build the Go binary
FROM golang:alpine AS builder
WORKDIR /app

# Copy the entire workspace context to resolve dependencies via go.work
COPY . .

# Build the specific server
RUN go build -o bin/mlc-markitdown mlc-markitdown/cmd/server/main.go

# Stage 2: Final runtime image
FROM python:3.11-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
    ffmpeg \
    libmagic1 \
    libxml2 \
    libxslt1-dev \
    gcc \
    python3-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
# Install markitdown
RUN pip install --no-cache-dir markitdown

# Copy Go binary and Python shim
COPY --from=builder /app/bin/mlc-markitdown /usr/local/bin/mlc-markitdown
COPY mlc-markitdown/internal/infrastructure/python/shim.py ./internal/infrastructure/python/shim.py

# Default artifact server address
# Note: When running on Linux host, you may need --add-host=host.docker.internal:host-gateway
ENV ARTIFACT_GRPC_ADDR=host.docker.internal:9590

# Entrypoint is the Go binary
ENTRYPOINT ["mlc-markitdown"]
