# Use a minimal Golang image to build the binary
FROM golang:1.16 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to download dependencies
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Use a smaller base image for the final stage
FROM debian:bullseye-slim

# Set the working directory inside the container
WORKDIR /app

# Copy only the binary from the builder stage
COPY --from=builder /app/main /app/main

# Install wkhtmltopdf dependencies
RUN apt-get update && apt-get install -y \
    wkhtmltopdf \
    && rm -rf /var/lib/apt/lists/*

EXPOSE 8080

# Command to run the executable
CMD ["./main"]