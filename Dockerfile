# Stage 1: Build the Go app
FROM golang:1.22.4-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire Go project into the container
COPY . .

# Build the Go binary for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/delivery-api

# Stage 2: Create the final image with a minimal footprint
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/delivery-api .

# Expose the port your app runs on
EXPOSE 8080

# Run the Go app
CMD ["./delivery-api"]
