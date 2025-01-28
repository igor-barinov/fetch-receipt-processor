# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy

# Copy the application source code and build the binary
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o fetch-server ./src/server

# Runtime stage
FROM alpine

# Copy only the binary from the build stage to the final image
COPY --from=builder /app/fetch-server /

# Expose the HTTP port
EXPOSE 3000

# Set the entry point for the container
ENTRYPOINT ["/fetch-server"]