# Use a minimal Go base image
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a minimal image for the final stage
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the built executable from the builder stage
COPY --from=builder /app/main .

# Change ownership of the executable
RUN chown appuser:appgroup /root/main

# Expose the port your application listens on (if applicable)
EXPOSE 8080

# Switch to the non-root user
USER appuser

# Run the executable
CMD ["./main"]