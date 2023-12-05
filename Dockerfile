# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Golang application files into the container
COPY . .

# Install wkhtmltopdf dependencies
RUN apt-get update && apt-get install -y \
    wget \
    fontconfig \
    libxrender1 \
    xfonts-75dpi \
    xfonts-base

# Download and install wkhtmltopdf
RUN wget https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6-1/wkhtmltox_0.12.6-1.bionic_amd64.deb

# Build the Golang application
RUN go build -o myapp

# Command to run the executable
CMD ["./myapp"]
