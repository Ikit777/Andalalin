# Use an official Golang runtime as a parent image
FROM golang:latest

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Install wkhtmltopdf dependencies
RUN apt-get update && apt-get install -y \
    libxrender1 \
    libfontconfig1 \
    libx11-dev \
    libjpeg62-turbo-dev \
    libpng-dev \
    libxext6 \
    fontconfig \
    wkhtmltopdf \
    && rm -rf /var/lib/apt/lists/*

# Build the Go app
RUN go mod download
RUN go build -o main .

EXPOSE 8080

# Command to run the executable
CMD ["./main"]