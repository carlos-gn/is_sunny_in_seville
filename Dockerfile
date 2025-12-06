# Build stage
FROM golang:1.25.5-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY main.go ./

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o seville-weather .

# Runtime stage - minimal image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /build/seville-weather .

# Run as non-root user
RUN adduser -D -u 1000 mcpuser
USER mcpuser

# Expose HTTP port
EXPOSE 8080

ENTRYPOINT ["/app/seville-weather"]
