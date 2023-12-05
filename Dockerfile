# Use an official Golang runtime as a parent image
FROM golang:latest

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Install wkhtmltopdf dependencies
RUN apt-get update && apt-get install -y \
    wkhtmltopdf \
    fontconfig \
    libfreetype6 \
    libjpeg62-turbo \
    libpng16-16 \
    libx11-6 \
    libxcb1 \
    libxext6 \
    libxrender1 \
    xfonts-75dpi \
    xfonts-base \
    && rm -rf /var/lib/apt/lists/*

# Build the Go app
RUN go build -o main .

EXPOSE 8080

# Command to run the executable
CMD ["./main"]