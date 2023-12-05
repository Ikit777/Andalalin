# Use the official Golang image as the base image
FROM golang:1.17-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy only the necessary files to the container
COPY go.mod go.sum ./

# Download and install Go module dependencies
RUN go mod download

# Install wkhtmltopdf dependencies
RUN apk --no-cache add \
    fontconfig \
    libfontconfig \
    libxrender \
    libxext \
    libintl \
    icu-libs \
    ttf-dejavu \
    && rm -rf /var/cache/apk/*

# Copy the rest of the application code to the container
COPY . .

# Build the Golang application
RUN go build -o main .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./main"]
