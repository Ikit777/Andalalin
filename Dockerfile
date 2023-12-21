# Use an official Golang runtime as a parent image
FROM golang:latest

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Update the package list
RUN apt-get update

# Install required dependencies
RUN apt-get install -y wkhtmltopdf

# Clean up unnecessary files
RUN rm -rf /var/lib/apt/lists/*

# Build the Go app
RUN go build -o main .

EXPOSE 8080

# Command to run the executable
CMD ["./main"]