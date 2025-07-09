# ---------- Stage 1: Build ----------
FROM golang:1.21-alpine AS builder

# Set environment variables
ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum separately to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary
RUN go build -o jalitalks main.go

# ---------- Stage 2: Final image ----------
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy only the built binary from builder
COPY --from=builder /app/jalitalks .

# Expose the app port (optional, if your app uses e.g., :8080)
EXPOSE 8080

# Run the Go binary
CMD ["./jalitalks"]
