# ---------------------------------------------------------
# Stage 1: Build the Go binary
# ---------------------------------------------------------
# We use the official GoCV image which contains the Go compiler 
# AND the exact matching version of OpenCV compiled from source.
FROM ghcr.io/hybridgroup/opencv:latest AS builder

WORKDIR /app

# Pre-copy and download go.mod dependencies to cache this layer
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
# -ldflags="-w -s" strips debugging information to make the binary much smaller
RUN go build -ldflags="-w -s" -o bin/homevision ./cmd/api


# ---------------------------------------------------------
# Stage 2: Create the runtime image
# ---------------------------------------------------------
# This image contains ONLY the OpenCV runtime libraries (no Go compiler),
# matching the exact OpenCV version used in the builder stage.
FROM ghcr.io/hybridgroup/opencv:latest

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/bin/homevision .

# Copy the static frontend files
COPY public/ ./public/

# Set production environment variables
ENV GIN_MODE=release
ENV PORT=8080

EXPOSE 8080

# Run the binary
CMD ["./homevision"]
