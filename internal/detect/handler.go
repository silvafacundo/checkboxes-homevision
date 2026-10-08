package detect

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DetectHandler handles image detection requests.
type DetectHandler struct {
	detector          Detector
	maxUploadSize     int64
	allowQueryOptions bool
}

// NewDetectHandler builds a DetectHandler.
func NewDetectHandler(d Detector, maxUploadSize int64, allowQueryOptions bool) *DetectHandler {
	return &DetectHandler{detector: d, maxUploadSize: maxUploadSize, allowQueryOptions: allowQueryOptions}
}

// Detect receives a multipart "image" field and runs detection on it.
func (h *DetectHandler) Detect(c *gin.Context) {
	// Limit request body size to prevent OOM
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, h.maxUploadSize)

	opts := DefaultDetectOptions()
	if h.allowQueryOptions {
		if err := c.ShouldBindQuery(&opts); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
			return
		}
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image provided"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open file"})
		return
	}
	defer file.Close()

	imgBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read file bytes"})
		return
	}

	result, err := h.detector.Detect(c.Request.Context(), imgBytes, opts)
	if err != nil {
		if errors.Is(err, ErrEmptyImage) || errors.Is(err, ErrInvalidImage) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Detection failed"})
		return
	}

	c.JSON(http.StatusOK, result)
}
