# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Install wkhtmltopdf dependencies
RUN apt-get update && apt-get install -y \
    wget \
    fontconfig \
    libxrender1 \
    xfonts-75dpi \
    xfonts-base

# Download and install wkhtmltopdf
RUN wget https://github.com/wkhtmltopdf/packaging/releases/0.12.6-1/wkhtmltox_0.12.6-1.bionic_amd64.deb
RUN dpkg -i wkhtmltox_0.12.6-1.bionic_amd64.deb

# Clean up
RUN rm wkhtmltox_0.12.6-1.bionic_amd64.deb

# Copy the Go source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

EXPOSE 8080

# Command to run the executable
CMD ["./main"]