# Use Go base image
FROM golang:1.23.2-bullseye

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum first, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the files
COPY . .

# Build the Go binary from the cmd/main.go file
RUN go build -o main ./cmd/main.go

# Expose port 8000
EXPOSE 8000

# Specify the command to run when the container starts
CMD ["/app/main"]