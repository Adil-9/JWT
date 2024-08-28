# Stage 1: Build the Go application
FROM golang:1.20-alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the app directory
COPY go.mod go.sum ./ 

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code from your working directory into the container
COPY . .

# Set the working directory to your cmd/medods-test-task directory
WORKDIR /app/cmd/medods-test-task

# Build the Go app and output the binary to /app/main
RUN go build -o /app/main

# Stage 2: Create a minimal image with the compiled binary
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=build /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
