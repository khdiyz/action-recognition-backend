package handler

import (
	"action-detector-backend/pkg/response"
	"fmt"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

const (
	maxFileSize = 50 << 20 // 50MB in bytes
)

// @Summary Upload a video file
// @Description Upload a video file with size limit of 50MB
// @Tags video
// @Accept multipart/form-data
// @Produce json
// @Param video formData file true "Video file to upload (max 50MB)"
// @Success 200 {object} response.BaseResponse "Returns the file link"
// @Failure 400 {object} response.BaseResponse "Invalid file size or type"
// @Failure 500 {object} response.BaseResponse "Internal server error"
// @Router /api/files [post]
func (h *Handler) uploadVideoFile(c *gin.Context) {
	// Get the file from the request
	file, err := c.FormFile("video")
	if err != nil {
		response.ErrorResponse(c, 400, err)
		return
	}

	// Check file size
	if file.Size > maxFileSize {
		response.ErrorResponse(c, 400, fmt.Errorf("file size exceeds maximum limit of 50MB"))
		return
	}

	// Check if file is a video
	if !isVideoFile(file) {
		response.ErrorResponse(c, 400, fmt.Errorf("file must be a video"))
		return
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		response.ErrorResponse(c, 500, err)
		return
	}
	defer src.Close()

	// Save the file or process it
	fileLink, err := h.usecase.UploadFile(c.Request.Context(), src, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		response.ErrorResponse(c, 500, err)
		return
	}

	// Return success response
	response.SuccessResponse(c, fileLink)
}

// isVideoFile checks if the uploaded file is a video based on its MIME type
func isVideoFile(file *multipart.FileHeader) bool {
	contentType := file.Header.Get("Content-Type")
	switch contentType {
	case "video/mp4",
		"video/quicktime",
		"video/x-msvideo",
		"video/x-matroska",
		"video/webm":
		return true
	default:
		return false
	}
}
