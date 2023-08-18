FROM ubuntu:latest

RUN apt-get update && \
    apt-get install -y wkhtmltopdf

# Build and run the Go application
CMD ["go", "run", "main.go"]