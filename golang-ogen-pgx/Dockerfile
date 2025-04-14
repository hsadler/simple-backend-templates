FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o app ./cmd/serverd

# Create a minimal production image
FROM alpine:3.18

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/app .
# COPY --from=builder /app/migrations ./migrations

# Set environment variables
ENV PORT=8000

# Expose the port
EXPOSE 8000

# Run the application
CMD ["./app"]
