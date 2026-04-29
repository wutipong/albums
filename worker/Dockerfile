# ---- Build Stage ----
FROM golang:1.26-alpine3.22 AS builder

# Install native dependencies required by govips and ffmpeg-go
RUN apk add --no-cache \
    pkgconfig \
    vips-dev \
    gcc \
    musl-dev

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o worker ./main.go

# ---- Runtime Stage ----
FROM alpine:3.22

# Install runtime native libraries (libvips, ffmpeg)
RUN apk add --no-cache \
    vips \
    vips-cpp \
    ffmpeg \
    ca-certificates

WORKDIR /app

COPY --from=builder /app/worker .

ENTRYPOINT ["/app/worker"]

