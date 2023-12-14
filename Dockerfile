# Use an official Golang image as a base image
FROM golang:1.16 as build

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Use a lightweight base image
FROM gcr.io/distroless/base

# Set the working directory inside the container
WORKDIR /

# Copy the binary from the build stage
COPY --from=build /app/golangapp /golangapp

# Install wkhtmltopdf
RUN apt-get update && apt-get install -y wkhtmltopdf

EXPOSE 8080

# Command to run the executable
CMD ["/golangapp"]