# 1. Build stage
FROM golang:1.24 AS builder

WORKDIR /app

# Copy modules first (enables caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/api

# 2. Final minimal runtime stage
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/server .
USER nonroot:nonroot
# Expose the port your app listens on (adjust if different)
EXPOSE 8080

# Run the server
CMD ["/app/server"]
