# Use the official Golang image as the base image
FROM golang:latest

# Install wkhtmltopdf dependencies
RUN apt-get update && apt-get install -y \
    wkhtmltopdf \
    && rm -rf /var/lib/apt/lists/*

# Set the working directory
WORKDIR /app

# Copy the Golang source code into the container
COPY . .

# Build the Golang application
RUN go build -o app

# Expose the port your application listens on
EXPOSE 8080

# Command to run the application
CMD ["./app"]
