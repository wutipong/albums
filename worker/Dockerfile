# ---- Build Stage ----
FROM golang:1.24-bookworm AS builder

# Install native dependencies required by govips and ffmpeg-go
RUN apt-get update && apt-get install -y --no-install-recommends \
    pkg-config \
    libvips-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o worker ./main.go

# ---- Runtime Stage ----
FROM debian:bookworm-slim

# Install runtime native libraries (libvips, ffmpeg)
RUN apt-get update && apt-get install -y --no-install-recommends \
    libvips42 \
    ffmpeg \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/worker .

ENTRYPOINT ["/app/worker"]

