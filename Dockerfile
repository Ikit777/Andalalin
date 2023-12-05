# Use an official Golang runtime as a parent image
FROM golang:latest

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Install wkhtmltopdf dependencies
RUN apt-get update && apt-get install -y \
    wget \
    fontconfig \
    libxrender1 \
    xfonts-75dpi \
    xfonts-base

# Download and install wkhtmltopdf
RUN wget https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6-1/wkhtmltox_0.12.6-1.bionic_amd64.deb
RUN dpkg -i wkhtmltox_0.12.6-1.bionic_amd64.deb

# Clean up
RUN rm wkhtmltox_0.12.6-1.bionic_amd64.deb

# Build the Go app
RUN go build -o main .

EXPOSE 8080

# Command to run the executable
CMD ["./main"]
