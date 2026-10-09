// package detect contains the application's business logic.
package detect

import (
	"context"
	"errors"
	"image"
	"math"

	"gocv.io/x/gocv"
)

var (
	// ErrEmptyImage is returned when an empty image is provided.
	ErrEmptyImage = errors.New("empty image")
	// ErrInvalidImage is returned when the bytes can't be decoded as an image.
	ErrInvalidImage = errors.New("invalid image")
)

// Checkbox is the internal bounding box of a detected checkbox, in pixels.
type Checkbox struct {
	X, Y, Width, Height int
}

// Box is a detected checkbox as returned by the API.
type Box struct {
	// BBox is [x1, y1, x2, y2] (top-left and bottom-right corners).
	BBox      [4]int `json:"bbox"`
	IsChecked bool   `json:"is_checked"`
}

// DetectionResult is the outcome of running detection on an image.
type DetectionResult struct {
	Boxes []Box `json:"boxes"`
}

// Detector runs detection over raw image bytes.
type Detector interface {
	Detect(ctx context.Context, img []byte, opts DetectOptions) (DetectionResult, error)
}

type detector struct{}

// NewDetector returns the default (OpenCV-based checkbox) Detector.
func NewDetector() Detector {
	return &detector{}
}

func (d *detector) Detect(_ context.Context, img []byte, opts DetectOptions) (DetectionResult, error) {
	if len(img) == 0 {
		return DetectionResult{}, ErrEmptyImage
	}

	// 1. Decode image and convert to grayscale.
	src, err := gocv.IMDecode(img, gocv.IMReadColor)
	if err != nil || src.Empty() {
		src.Close()
		return DetectionResult{}, ErrInvalidImage
	}
	defer src.Close()

	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(src, &gray, gocv.ColorBGRToGray)

	// 2. Binarize (inverted): background black, lines white.
	thresh := gocv.NewMat()
	defer thresh.Close()

	gocv.AdaptiveThreshold(gray, &thresh, opts.ThresholdMaxVal,
		gocv.AdaptiveThresholdGaussian, gocv.ThresholdBinaryInv,
		opts.AdaptiveBlockSize, opts.AdaptiveC)

	// 3. Find all contours.
	contours := gocv.FindContours(thresh, gocv.RetrievalList, gocv.ChainApproxSimple)
	defer contours.Close()

	var candidates []Checkbox

	// 4. Keep only square-ish 4-vertex contours within the area range.
	for i := 0; i < contours.Size(); i++ {
		cnt := contours.At(i)

		perimeter := gocv.ArcLength(cnt, true)

		approx := gocv.ApproxPolyDP(cnt, opts.ApproxEpsilonFrac*perimeter, true)
		vertices := approx.Size()

		var rect image.Rectangle
		if vertices == opts.VerticesCount {
			rect = gocv.BoundingRect(approx)
		}
		approx.Close()

		// Check if the shape has 4 vertices (VerticesCount is typically 4 for a rectangle)
		if vertices != opts.VerticesCount || rect.Dy() == 0 {
			continue
		}

		aspectRatio := float64(rect.Dx()) / float64(rect.Dy())
		area := gocv.ContourArea(cnt)

		// Check if it has a square-like aspect ratio and a reasonable size for a checkbox
		if aspectRatio >= opts.MinAspectRatio && aspectRatio <= opts.MaxAspectRatio &&
			area > opts.MinArea && area < opts.MaxArea {
			box := Checkbox{
				X: rect.Min.X, Y: rect.Min.Y, Width: rect.Dx(), Height: rect.Dy(),
			}
			candidates = append(candidates, box)
		}
	}

	// 5. Remove overlapping detections (e.g. inner and outer contour of the same box).
	checkboxes := removeCollisions(candidates, opts.CollisionRelThres)

	boxes := make([]Box, 0, len(checkboxes))
	for _, cb := range checkboxes {
		boxes = append(boxes, Box{
			BBox:      [4]int{cb.X, cb.Y, cb.X + cb.Width, cb.Y + cb.Height},
			IsChecked: isChecked(thresh, cb, opts),
		})
	}
	return DetectionResult{Boxes: boxes}, nil
}

// removeCollisions drops boxes whose center is closer than
// relativeThreshold * width of an already accepted box.
func removeCollisions(boxes []Checkbox, relativeThreshold float64) []Checkbox {
	final := make([]Checkbox, 0, len(boxes))
	for _, b := range boxes {
		cx := float64(b.X) + float64(b.Width)/2
		cy := float64(b.Y) + float64(b.Height)/2

		duplicated := false
		for _, f := range final {
			fcx := float64(f.X) + float64(f.Width)/2
			fcy := float64(f.Y) + float64(f.Height)/2

			if math.Hypot(cx-fcx, cy-fcy) < float64(f.Width)*relativeThreshold {
				duplicated = true
				break
			}
		}
		if !duplicated {
			final = append(final, b)
		}
	}
	return final
}

// isChecked determines if a box is checked by calculating the percentage
// of white pixels inside a shrunken inner bounding box.
func isChecked(src gocv.Mat, checkbox Checkbox, opts DetectOptions) bool {
	marginX := int(float64(checkbox.Width) * opts.InnerMarginFracX)
	marginY := int(float64(checkbox.Height) * opts.InnerMarginFracY)

	innerRect := image.Rect(
		checkbox.X+marginX,
		checkbox.Y+marginY,
		checkbox.X+checkbox.Width-marginX,
		checkbox.Y+checkbox.Height-marginY,
	)

	roi := src.Region(innerRect)

	whitePixels := gocv.CountNonZero(roi)
	roi.Close()

	totalArea := float64(innerRect.Dx() * innerRect.Dy())
	filled := (float64(whitePixels) / totalArea) * 100

	return filled > opts.IsCheckedThres
}
