# Build stage
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . ./

# Build the application binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /userapi ./cmd/api

# Run stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /

# Copy the built binary from the builder stage
COPY --from=builder /userapi /userapi

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["/userapi"]