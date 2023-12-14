# Use an official Golang runtime as a parent image
FROM golang:latest

# Set the working directory to /app
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the current directory contents into the container at /app
COPY . .

# Install wkhtmltopdf dependencies
RUN apt-get update && apt-get install -y \
    wkhtmltopdf \
    && rm -rf /var/lib/apt/lists/*

# Build the Go app
RUN go build -o main .

# Add permissions for files
RUN chmod +x main
RUN chmod -R +r 644 /app/templates

EXPOSE 8080

# Command to run the executable
CMD ["./main"]