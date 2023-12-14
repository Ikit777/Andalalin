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
    wget \
    unzip \
    net-tools \
    vim \
    && rm -rf /var/lib/apt/lists/*

# Download and install wkhtmltopdf separately
RUN wget https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6-1/wkhtmltox_0.12.6-1.bionic_amd64.deb && \
    dpkg -i wkhtmltox_0.12.6-1.bionic_amd64.deb && \
    apt install -f

# Build the Go app
RUN go mod download
RUN go build -o main .

EXPOSE 8080

# Command to run the executable
CMD ["./main"]