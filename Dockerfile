# Use an official Golang runtime as a parent image
FROM golang:latest

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Install wkhtmltopdf dependencies
RUN apt-get update && apt-get install -y \
    curl \
    libxrender1 \
    libjpeg62-turbo \
    fontconfig \
    libxtst6 \
    xfonts-75dpi \
    xfonts-base \
    xz-utils \
    && rm -rf /var/lib/apt/lists/*
    
RUN curl "https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6-1/wkhtmltox_0.12.6-1.buster_amd64.deb" -L -o "wkhtmltopdf.deb"
RUN dpkg -i wkhtmltopdf.deb

# Build the Go app
RUN go build -o main .

EXPOSE 8080

# Command to run the executable
CMD ["./main"]