# Use a smaller base image for production
FROM golang:alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the local code to the container
COPY . .

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

# Build the Golang application
RUN go build -o main .

# Use a minimal base image for the final image
FROM alpine:latest

# Copy only necessary files from the build image
COPY --from=build /app/main /app/main

# Expose the port the app runs on
EXPOSE 8080

# Run the application
CMD ["/app/main"]
