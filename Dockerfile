# Use the official golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the local code to the container
COPY . .

# Install wkhtmltopdf dependencies
RUN apt-get update && \ apt-get install -y \ fontconfig \ libfontconfig1 \ libfreetype6 \ libx11-6 \ libxext6 \ libxrender1 \ xfonts-base \ xfonts-75dpi \ wget \ && apt-get clean \ && rm -rf /var/lib/apt/lists/*

# Download and install wkhtmltopdf
RUN wget https://github.com/wkhtmltopdf/wkhtmltopdf/releases/download/0.12.6/wkhtmltox_0.12.6-1.bionic_amd64.deb && \ dpkg -i wkhtmltox_0.12.6-1.bionic_amd64.deb && \ apt-get install -f

# Build the Golang application
RUN go build -o main .

# Expose the port the app runs on
EXPOSE 8080

# Run the application
CMD ["./main"]
