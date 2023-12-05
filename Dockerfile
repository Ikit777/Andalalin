# Use an official Golang runtime as a parent image
FROM golang:latest AS builder

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Build the Go app
RUN go build -o main .

# Stage 2: Create a lightweight image
FROM debian:bullseye-slim

# Install wkhtmltopdf dependencies and any other necessary dependencies
RUN apt-get update && apt-get install -y \
    wkhtmltopdf \
    && rm -rf /var/lib/apt/lists/*

# Set the working directory to /app
WORKDIR /app

# Copy the built Go binary from the builder stage
COPY --from=builder /app/main .

# Expose the port on which your Golang app runs
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
