package detect

type DetectOptions struct {
	// --- 1. Thresholding (Binarization) ---

	// Pixel value assigned to areas that pass the threshold (usually 255 for pure white).
	ThresholdMaxVal float32 `form:"threshold_max_val"`
	// Size of the pixel neighborhood used to calculate the local threshold. Must be an odd number.
	AdaptiveBlockSize int `form:"adaptive_block_size"`
	// Constant subtracted from the mean to fine-tune adaptive thresholding and remove noise.
	AdaptiveC float32 `form:"adaptive_c"`

	// --- 2. Contour Approximation ---

	// Multiplier for contour perimeter to determine polygon approximation strictness.
	ApproxEpsilonFrac float64 `form:"approx_epsilon_frac"`
	// Required number of polygon vertices a contour must have (typically 4 for a rectangle).
	VerticesCount int `form:"vertices_count"`

	// --- 3. Geometric Filters ---

	// Minimum width-to-height ratio allowed for a bounding box.
	MinAspectRatio float64 `form:"min_aspect_ratio"`
	// Maximum width-to-height ratio allowed for a bounding box.
	MaxAspectRatio float64 `form:"max_aspect_ratio"`
	// Minimum pixel area required for a contour to be considered a checkbox.
	MinArea float64 `form:"min_area"`
	// Maximum pixel area allowed for a contour.
	MaxArea float64 `form:"max_area"`

	// --- 4. Collision Detection ---

	// Distance multiplier used to eliminate overlapping duplicate detections.
	CollisionRelThres float64 `form:"collision_rel_thres"`

	// --- 5. Checked State ---

	// Percentage of the width to shrink inwards when checking if the box is filled.
	InnerMarginFracX float64 `form:"inner_margin_frac_x"`
	// Percentage of the height to shrink inwards when checking if the box is filled.
	InnerMarginFracY float64 `form:"inner_margin_frac_y"`
	// Minimum percentage of filled pixels inside the inner rectangle to classify as "checked".
	IsCheckedThres float64 `form:"is_checked_thres"`
}

func DefaultDetectOptions() DetectOptions {
	return DetectOptions{
		// 1. Thresholding
		ThresholdMaxVal:   255.0,
		AdaptiveBlockSize: 15,
		AdaptiveC:         5.0,

		// 2. Contour Approximation
		ApproxEpsilonFrac: 0.02,
		VerticesCount:     4,

		// 3. Geometric Filters
		MinAspectRatio: 0.85,
		MaxAspectRatio: 1.3,
		MinArea:        200.0,
		MaxArea:        5000.0,

		// 4. Collision Detection
		CollisionRelThres: 0.4,

		// 5. Checked State
		InnerMarginFracX: 0.15,
		InnerMarginFracY: 0.15,
		IsCheckedThres:   10.0,
	}
}
