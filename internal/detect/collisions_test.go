package detect

import (
	"testing"
)

func TestRemoveCollisions(t *testing.T) {
	tests := []struct {
		name              string
		boxes             []Checkbox
		relativeThreshold float64
		expectedCount     int
	}{
		{
			name: "No collisions",
			boxes: []Checkbox{
				{X: 10, Y: 10, Width: 100, Height: 100},
				{X: 200, Y: 200, Width: 100, Height: 100}, // Far away
			},
			relativeThreshold: 0.4,
			expectedCount:     2, // Both should be kept
		},
		{
			name: "Exact duplicate",
			boxes: []Checkbox{
				{X: 10, Y: 10, Width: 100, Height: 100},
				{X: 10, Y: 10, Width: 100, Height: 100},
			},
			relativeThreshold: 0.4,
			expectedCount:     1, // One should be dropped
		},
		{
			name: "Slightly shifted overlapping box (collision)",
			boxes: []Checkbox{
				{X: 10, Y: 10, Width: 100, Height: 100},
				{X: 15, Y: 15, Width: 100, Height: 100}, // Center shifted by just 5 pixels
			},
			relativeThreshold: 0.4,
			// Threshold is 40 pixels (100 * 0.4). Hypotenuse of (5,5) is ~7.07, so it's a collision.
			expectedCount: 1, 
		},
		{
			name: "Multiple collisions and one distinct",
			boxes: []Checkbox{
				{X: 10, Y: 10, Width: 100, Height: 100},
				{X: 12, Y: 12, Width: 100, Height: 100},   // Collides with 1st
				{X: 8, Y: 9, Width: 100, Height: 100},     // Collides with 1st
				{X: 500, Y: 500, Width: 100, Height: 100}, // Distinct, far away
			},
			relativeThreshold: 0.4,
			expectedCount:     2, // Should keep the 1st and the 4th
		},
		{
			name: "Barely outside relative threshold (no collision)",
			boxes: []Checkbox{
				{X: 10, Y: 10, Width: 100, Height: 100},
				{X: 50, Y: 10, Width: 100, Height: 100}, // Shifted by 40 pixels horizontally
			},
			// Threshold is 30 pixels (100 * 0.3). Distance is 40, so NO collision.
			relativeThreshold: 0.3,
			expectedCount:     2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeCollisions(tt.boxes, tt.relativeThreshold)
			if len(result) != tt.expectedCount {
				t.Errorf("removeCollisions() returned %d boxes, want %d", len(result), tt.expectedCount)
			}
		})
	}
}
