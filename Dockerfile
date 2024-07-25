# Use the official golang image to create a build artifact
FROM golang:1.22-alpine3.20 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod file
COPY go.mod ./

# Download all dependencies.
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o httprelay .

# Start a new stage from scratch
FROM alpine:3.20

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/httprelay /httprelay

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["/httprelay"]