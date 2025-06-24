# Stage 1: Build the Go application
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
# CGO_ENABLED=0 for static linking, GOOS=linux for cross-compilation if needed
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o proxypanel-api ./cmd/proxypanel-api/main.go

# Stage 2: Create a lightweight image
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/proxypanel-api .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./proxypanel-api"]
