# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

# Install gcc and dependencies for CGO
RUN apk add --no-cache gcc musl-dev  # <-- THIS FIXES THE ERROR

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o /snippet-manager

# Stage 2: Create the final lightweight image
FROM alpine:latest

# Install SQLite (required for your database)
RUN apk --no-cache add sqlite



# Copy the binary from builder
COPY --from=builder /snippet-manager /snippet-manager

# Copy existing database file (if any)
# COPY database/snippets.db /database/snippets.db

# Set permissions (important for SQLite in Docker)
RUN chmod +x /snippet-manager

# Expose port (match your Gin app's port)
EXPOSE 8080

# Command to run the application
CMD ["/snippet-manager"]