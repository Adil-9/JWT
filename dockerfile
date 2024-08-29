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