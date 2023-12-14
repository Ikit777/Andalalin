# Use an official Golang runtime as a parent image
FROM golang:latest

# Set the working directory to /app
WORKDIR /app

# Install wkhtmltopdf dependencies
RUN apt-get update && \
    apt-get install -y libxrender1 libfontconfig1 libx11-dev libxext-dev libfreetype6 libjpeg62-turbo libpng16-16

# Install go-wkhtmltopdf
RUN go get -u github.com/SebastiaanKlippert/go-wkhtmltopdf

# Copy the local package files to the container's workspace
COPY . /app

# Build the Go app
RUN go build -o main .

EXPOSE 8080

# Command to run the executable
CMD ["./main"]