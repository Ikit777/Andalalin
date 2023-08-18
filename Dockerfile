# Use the official Golang image as the base
FROM golang:latest

# Set the working directory
WORKDIR /app

# Copy the wkhtmltopdf binaries into the container
COPY wkhtmltopdf_bin /app/wkhtmltopdf_bin

# Copy the Go application files to the container
COPY . .

# Build and run the Go application
CMD ["go", "run", "main.go"]