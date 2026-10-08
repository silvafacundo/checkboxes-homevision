# Homevision

Gin-based HTTP API and web interface for OpenCV checkbox detection.

## Project Structure

```text
.
├── cmd/
│   └── api/            # Application entrypoint (main package)
├── internal/
│   └── detect/         # Detection services, config, handlers, and router
├── public/             # Static frontend files (index.html, styles)
├── Dockerfile          # Multi-stage Docker build
├── .env                # Environment configuration variables
├── go.mod
└── go.sum
```

## Setup & Running

You have two options for running this project: using Docker (Recommended) or running it locally via Go. 

### Option 1: Via Docker (Recommended)

Running via Docker is highly recommended because compiling OpenCV and GoCV locally can be complicated, particularly on Windows. The provided Dockerfile uses the official GoCV images.

1. **Build the image**:
   ```sh
   docker build -t homevision .
   ```

2. **Run the container**:
   ```sh
   docker run -p 8080:8080 --name homevision-app homevision
   ```

3. **Access the application**: Open your browser and navigate to `http://localhost:8080`.

### Option 2: Local Setup (Requires OpenCV toolchain)

To run the project locally, you must have the Go compiler, GCC/MinGW (on Windows), and the OpenCV C++ libraries installed. 

**Prerequisites:**
- [Go](https://golang.org/dl/)
- OpenCV 4.x compiled and installed locally (see [GoCV Windows installation guide](https://gocv.io/getting-started/windows/)).

**Running:**
Run the application:
   ```sh
   go run ./cmd/api
   ```

**Testing:**
```sh
go test ./...
```
*(Note: To run tests locally, your OpenCV toolchain must be properly configured so CGO can compile the OpenCV bindings).*

## Configuration

You can configure the application using a `.env` file or environment variables. 

| Variable               | Default | Description                                             |
|------------------------|---------|---------------------------------------------------------|
| `PORT`                 | `8080`  | HTTP port to bind to                                    |
| `GIN_MODE`             | `debug` | Gin mode (`debug`, `release`, `test`)                   |
| `MAX_UPLOAD_SIZE_MB`   | `20`    | Maximum allowed file upload size (MB)                   |
| `ALLOW_QUERY_OPTIONS`  | `false` | Allow overriding detector options via query params      |

## API Endpoints

- `GET /` - Serves the web interface
- `GET /api/ping` - Health check
- `POST /api/detect` - Upload an image to detect checkboxes

**Example API Request:**
```sh
curl -F "image=@photo.jpg" http://localhost:8080/api/detect
```

## Known Limitations

- **Tilted Checkboxes**: The current contour detection algorithm expects boxes to be axis-aligned. Checkboxes that are significantly rotated or skewed are not supported.
- **Image Noise**: Images with heavy artifacting, noise, or poor lighting may degrade detection accuracy and lead to suboptimal results. For best performance, use clear, well-lit scans or high-quality photos.
- **Concurrency**: This service performs CPU-intensive image processing synchronously on the HTTP request thread. It is not currently optimized for high concurrency or heavy production traffic. To support a large volume of requests, a dedicated background queue/worker architecture (e.g., using Redis or RabbitMQ) would be required to process images asynchronously.
