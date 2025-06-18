# 1. Use an official Go image as a builder
FROM golang:1.21 as builder

# 2. Set working directory inside the container
WORKDIR /app

# 3. Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# 4. Copy the source code into the container
COPY . .

# 5. Build the Go app
RUN go build -o main .

# 6. Use a small base image for final container
FROM debian:bullseye-slim

# 7. Set working directory again
WORKDIR /app

# 8. Copy the binary from the builder
COPY --from=builder /app/main .

# 9. Expose port (change it to match your app)
EXPOSE 8080

# 10. Command to run the app
CMD ["./main"]
