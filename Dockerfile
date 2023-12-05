# Use an official Golang runtime as a parent image
FROM golang:latest

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Copy the built Go binary from the builder stage
COPY --from=builder /app/main .

# Install wkhtmltopdf dependencies and any other necessary dependencies
RUN apt-get update && apt-get install -y \
    wkhtmltopdf \
    && rm -rf /var/lib/apt/lists/*

# Build the Go app
RUN go build -o main .

# Expose the port on which your Golang app runs
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
