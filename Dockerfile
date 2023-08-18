# Use the official Golang image as the base
FROM golang:latest

# Install wkhtmltopdf dependencies and wkhtmltopdf itself
RUN apt-get update && apt-get install -y \
    libxrender1 \
    libfontconfig1 \
    libx11-dev \
    libjpeg62-turbo-dev \
    libxext6 \
    wkhtmltopdf
    
# Set the working directory
WORKDIR /app

# Copy the Go application files to the container
COPY . .

# Build and run the Go application
CMD ["go", "run", "main.go"]