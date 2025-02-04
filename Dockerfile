# Use the official Golang image as the base image
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application code to the working directory
COPY . .

# Build the Go application
RUN go build -o receipt-processor main.go

# Expose the port the application will run on
EXPOSE 8080

# Command to run the application
CMD ["./receipt-processor"]