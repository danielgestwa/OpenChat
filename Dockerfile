# --- Build Stage ---
FROM golang:1.25 AS builder

# Set working directory inside container
WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Replace env variables with prod setup
COPY .env_prod .env

# Build the Go app (static binary for smaller image)
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# --- Run Stage ---
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Copy templates, static files, images, etc.
COPY templates/ ./templates/
COPY lib/ ./lib/
COPY images/ ./images/

# Expose HTTP port
EXPOSE 80

# Run the application
CMD ["./server"]
