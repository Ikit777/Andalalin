# Use an official Golang runtime as a parent image
FROM frolvlad/alpine-glibc:alpine-3.14 AS builder

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Build the Go app
RUN go build -o main .

# Stage 2: Create a lightweight image
FROM alpine:latest

# Set the working directory to /app
WORKDIR /app

# Install wkhtmltopdf dependencies and any other necessary dependencies
RUN apk --no-cache add \
    wkhtmltopdf

# Copy the built Go binary from the builder stage
COPY --from=builder /app/main .

EXPOSE 8080

# Command to run the executable
CMD ["./main"]
