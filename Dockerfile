# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies and create non-root user
RUN apk add --no-cache gcc musl-dev git && \
    adduser -D -u 10001 appuser

# Set working directory
WORKDIR /app

# Copy go mod and sum files first (for better caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy the rest of the application code
COPY . .

# Install Swagger CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger docs
RUN swag init -g main.go --parseDependency --parseInternal --parseDepth 1

# Build the application with optimization
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a \
    -installsuffix cgo \
    -ldflags="-w -s" \
    -o main .

# Final stage - minimal runtime image
FROM alpine:latest

# Install CA certificates and create non-root user
RUN apk --no-cache add ca-certificates && \
    adduser -D -u 10001 appuser

# Set working directory
WORKDIR /app

# Copy binary and necessary files from builder
COPY --from=builder --chown=appuser:appuser /app/main .
COPY --from=builder --chown=appuser:appuser /app/config ./config
COPY --from=builder --chown=appuser:appuser /app/docs ./docs

# Set environment variables with secure defaults
ENV GIN_MODE=release \
    MONGODB_URI=mongodb://mongodb:27017 \
    DB_NAME=taskify \
    SERVER_ADDRESS=0.0.0.0 \
    SERVER_PORT=3000

# Expose application port
EXPOSE 3000

# Switch to non-root user for security
USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s \
    CMD wget --spider http://localhost:3000/health || exit 1

# Run the application
CMD ["./main"]
