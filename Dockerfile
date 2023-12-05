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
    xfonts-base \
    libjpeg62-turbo \
    libx11-6 \
    libxcb1 \
    libxext6 \
    libxrender1 \
    libssl1.1 \
    libx11-xcb1

# Download and install wkhtmltopdf
RUN wget https://github.com/wkhtmltopdf/wkhtmltopdf/releases/download/LATEST_VERSION/wkhtmltox_LATEST_VERSION.bionic_amd64.deb
RUN dpkg -i wkhtmltox_LATEST_VERSION.bionic_amd64.deb

# Build the Go app
RUN go build -o main .

EXPOSE 8080

# Command to run the executable
CMD ["./main"]
