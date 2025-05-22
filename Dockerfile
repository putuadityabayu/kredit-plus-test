# Stage 1: Build the application
FROM golang:1.23.6-alpine3.21 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to leverage Docker's layer caching
COPY go.mod go.sum ./

# Download Go modules (dependencies)
# Using --mount=type=cache for go module cache
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
# Using --mount=type=cache for Go's build cache
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o /app/main -ldflags="-s -w" main.go

# Stage 2: Create the final, minimal production image
FROM gcr.io/distroless/static-debian11

# Set the working directory inside the container
WORKDIR /app

# Copy the compiled binary from the 'builder' stage to the final image
COPY --from=builder /app/main .

# Expose the port your application will listen on
EXPOSE 4000

# Define the command to run your application when the container starts
CMD ["./main", "server"]