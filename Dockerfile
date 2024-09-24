# Stage 1: Build the Go application
FROM golang:1.22.0 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first for dependency resolution
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api ./cmd/api

# Stage 2: Build a small image with only the built binary
#FROM gcr.io/distroless/base-debian12

FROM debian:bookworm-slim
# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/api .

RUN apt-get update && apt-get install -y \
    bash \
    && rm -rf /var/lib/apt/lists/*

# Expose the application port (if the app listens on a specific port)
EXPOSE 8080
ENTRYPOINT ["/bin/bash"]

# Command to run the application
CMD ["/app/api"]