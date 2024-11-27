# Use the latest Golang image as the base image
FROM golang:1.23.3 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod verify
RUN go mod tidy

# Copy the rest of the application source code
COPY . .

# Build the application
RUN go build -o ./build/app ./cmd/app/main.go

# Expose the port the application listens on
EXPOSE 3000

# Run the application
CMD ["sh", "-c", "./build/app"]