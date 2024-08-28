# # Stage 1: Build the Go application
# FROM golang:1.23-alpine AS build

# # Set the Current Working Directory inside the container
# WORKDIR /app

# # Copy go.mod and go.sum files to the app directory
# COPY go.mod go.sum ./ 

# # Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
# RUN go mod download

# # Copy the source code from your working directory into the container
# COPY . .

# # Set the working directory to your cmd/medods-test-task directory
# WORKDIR /app/cmd/medods-test-task

# # Build the Go app and output the binary to /app/main
# RUN go build -o /app/main
# # go build -o /app/main /cmd/medods-test-task

# # Stage 2: Create a minimal image with the compiled binary
# FROM alpine:latest

# # Set the Current Working Directory inside the container
# WORKDIR /root/

# ENV secret_key="hello-world!"
# ENV host="localhost"
# ENV port="5432"
# ENV user="postgres"
# ENV password="postgres"
# ENV dbname="postgres"

# # Copy the Pre-built binary file from the previous stage
# COPY --from=build /app/main .

# # Expose port 8080 to the outside world
# EXPOSE 8080
# EXPOSE 5432

# # Command to run the executable
# CMD ["./main"]


# Use an official Go runtime as a parent image
FROM golang:alpine

# Set the working directory in the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Build the Go application
RUN go get -d -v ./... && \
    go install -v ./...

ENV secret_key="hello-world!"
ENV host="localhost"
ENV port="5432"
ENV user="postgres"
ENV password="postgres"
ENV dbname="postgres"

CMD ["go", "run", "cmd/medods-test-task/main.go"]