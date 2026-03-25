# --- Stage 1: Build ---
FROM golang:1.26-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Leverage Docker cache for dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary with optimizations (omit symbol table and DWARF)
RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-s -w" \
  -o /bin/api ./cmd/api

# --- Stage 2: Runtime ---
FROM alpine:3.23

# Add standard libraries and timezome data
RUN apk add --no-cache ca-certificates tzdata

# Create a non-privileged user for security
RUN adduser -D -g '' appuser
USER appuser

COPY --from=builder /bin/api /bin/api

# Standard API port
EXPOSE 4000

# Start the application
ENTRYPOINT [ "/bin/api" ]
