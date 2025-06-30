# Use the official Go image as the base image
FROM golang:1.24.4 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first for better caching of dependencies
COPY go.mod ./
# Generate go.sum and download dependencies
RUN go mod download && go mod tidy

# Copy the rest of the source code into the container
COPY . .

# Build the application with a target to the main.go file
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/api/main.go

# Use a minimal Alpine image for the final stage
FROM alpine:latest
WORKDIR /root/

# Install ca-certificates, important for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the executable from the builder stage
COPY --from=builder /app/server .
COPY --from=builder /app/.env .

# Expose the port the application will run on
EXPOSE 8080

# Run the compiled server binary
CMD ["./server"]
