package detect_test

import (
	"context"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"testing"

	"homevision/internal/detect"
)

func TestDetector(t *testing.T) {
	// Find all .json files in the testdata directory as the source of truth for test cases
	jsonFiles, err := filepath.Glob("testdata/*.json")
	if err != nil {
		t.Fatalf("Failed to glob testdata: %v", err)
	}
	if len(jsonFiles) == 0 {
		t.Fatalf("No test JSON files found in testdata")
	}

	for _, jsonPath := range jsonFiles {
		// Strip the .json extension to find the base name
		baseName := jsonPath[:len(jsonPath)-5]

		// Look for a corresponding image file
		var imgPath string
		for _, ext := range []string{".png", ".jpg", ".jpeg", ".webp"} {
			if _, err := os.Stat(baseName + ext); err == nil {
				imgPath = baseName + ext
				break
			}
		}

		if imgPath == "" {
			t.Errorf("Missing image file for test case: %s", jsonPath)
			continue
		}

		t.Run(filepath.Base(imgPath), func(t *testing.T) {
			// Read test image
			imgBytes, err := os.ReadFile(imgPath)
			if err != nil {
				t.Fatalf("Failed to read image %q: %v", imgPath, err)
			}

			// Read corresponding expected JSON output
			jsonBytes, err := os.ReadFile(jsonPath)
			if err != nil {
				t.Fatalf("Failed to read expected json %q: %v", jsonPath, err)
			}

			var expected detect.DetectionResult
			if err := json.Unmarshal(jsonBytes, &expected); err != nil {
				t.Fatalf("Failed to parse expected json: %v", err)
			}

			// Initialize detector and run
			detector := detect.NewDetector()
			opts := detect.DefaultDetectOptions()

			result, err := detector.Detect(context.Background(), imgBytes, opts)
			if err != nil {
				t.Fatalf("Detect returned error: %v", err)
			}

			// Compare using flexible overlap logic
			matchBoxes(t, expected.Boxes, result.Boxes)
		})
	}
}

// matchBoxes checks that every expected box has a corresponding detected box
// that physically overlaps it and shares the same is_checked value.
func matchBoxes(t *testing.T, expected, actual []detect.Box) {
	if len(expected) != len(actual) {
		t.Errorf("Expected %d boxes, got %d", len(expected), len(actual))
		return
	}

	matched := make([]bool, len(actual))

	for _, e := range expected {
		ecx := float64(e.BBox[0]+e.BBox[2]) / 2.0
		ecy := float64(e.BBox[1]+e.BBox[3]) / 2.0
		ew := float64(e.BBox[2] - e.BBox[0])

		foundMatch := false
		for j, a := range actual {
			if matched[j] {
				continue
			}
			acx := float64(a.BBox[0]+a.BBox[2]) / 2.0
			acy := float64(a.BBox[1]+a.BBox[3]) / 2.0

			// If centers are within 50% of the box width, they are the same box
			if math.Hypot(ecx-acx, ecy-acy) < ew*0.5 {
				if e.IsChecked != a.IsChecked {
					t.Errorf("Box at (%d, %d) matched, but IsChecked differs. Expected %v, got %v", 
						int(ecx), int(ecy), e.IsChecked, a.IsChecked)
				}
				matched[j] = true
				foundMatch = true
				break
			}
		}

		if !foundMatch {
			t.Errorf("Could not find a matching overlapping box for expected box at %v", e.BBox)
		}
	}
}
