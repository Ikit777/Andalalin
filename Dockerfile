# Use a Golang base image
FROM golang:lates

# Install wkhtmltopdf dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    wkhtmltopdf \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy your Golang application code into the container
COPY . .

# Build your Golang application
RUN go build -o app

# Start your Golang application
CMD ["./app"]