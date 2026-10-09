# Homevision

Gin-based HTTP API and web interface for OpenCV checkbox detection.

## Table of Contents
- [Live Demo](#live-demo)
- [Setup & Running](#setup--running)
- [Configuration](#configuration)
- [API Endpoints](#api-endpoints)
- [Engineering Notes](#engineering-notes)

## Live Demo

A live demo of the application with a web interface is deployed and available at:
**[https://checkbox-detector.facus.ar/](https://checkbox-detector.facus.ar/)**

I chose to deploy this service on Render because it natively supports deploying directly from a `Dockerfile`, which made it incredibly straightforward to host the necessary OpenCV dependencies.

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

**Example Response:**
```json
{
  "boxes": [
    {
      "bbox": [150, 300, 180, 330],
      "is_checked": true
    },
    {
      "bbox": [150, 360, 180, 390],
      "is_checked": false
    }
  ]
}
```

## Engineering Notes

**Personal Note:** I used this challenge as an opportunity to step out of my comfort zone and learn two technologies I hadn't worked with before: **Go** and **Computer Vision**.

### Design Decisions & Approach
When I looked at the example images provided for this challenge, I assumed that the system would generally process clear images with square checkboxes. Because of this assumption, I decided to use a direct computer vision approach (OpenCV) rather than trying to train a Machine Learning model. This decision made sense to me because it is a fast and simple solution that easily handles the provided examples, and it completely avoids the need to gather and label a large training dataset.

The core algorithm uses OpenCV and works in three main steps:
1. **Preprocessing:** The image is converted to grayscale, and an adaptive threshold is applied to obtain a binary (black and white) image.
2. **Contour Detection:** OpenCV's contour algorithm finds shapes in the binary image, which the program filters to identify squares that could be checkboxes.
3. **State Evaluation:** Once a checkbox is detected, I calculate the percentage of non-background pixels inside it. If this percentage is greater than a defined threshold, the checkbox is considered "checked".

### Limitations & Downsides
While this approach works well for the provided examples, relying on those assumptions leads to a few weak points:
- **Image Quality:** Because I assumed the inputs would be clear and square, the results can be heavily affected by image noise, poor quality, or if the checkboxes are marked heavily over the edges.
- **Rotation:** It does not support checkboxes or images rotated at an angle. While the contour algorithm can identify rotated shapes, the logic that evaluates if the box is checked currently assumes a straight, axis-aligned square.
- **Concurrency:** This service performs CPU-intensive image processing synchronously on the HTTP request thread. It would struggle with high concurrency in its current state.

### Future Improvements
If I were to prepare this for a high-traffic production environment, I would opt for an asynchronous queue/worker architecture (e.g., using Redis or RabbitMQ) to process the images without blocking the web server.