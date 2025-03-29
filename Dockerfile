FROM golang:1.22.2
LABEL maintainer="Moses Onyango, Aaron Ochieng, Swabri Musa, Andy Osindo"
LABEL version="1.0"
LABEL description="An advanced golang web project"

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY ./assets ./assets
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./Makefile ./

# Build the application
RUN go build -o forum ./cmd/web/

# Expose the port
EXPOSE 8000

# Run the application
CMD ["./forum"]