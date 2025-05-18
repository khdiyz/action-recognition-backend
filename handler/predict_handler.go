package handler

import (
	"action-detector-backend/models"
	"action-detector-backend/pkg/response"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// pretictActionResponse represents the response structure for action prediction
// @Description Response containing predicted actions in Uzbek language
type pretictActionResponse struct {
	// @Description List of predicted actions in Uzbek language
	Predictions []models.Prediction `json:"predictions"`
	VideoURL    string              `json:"output_video_url"`
}

// predictActionRequest represents the request body for predicting actions
// @Description Request containing the URL of the video to analyze
type predictActionRequest struct {
	// @Description URL of the video to analyze
	VideoURL      string `json:"video_url" binding:"required"`
	PredictApiURL string `json:"predict_api_url"`
}

// @Summary Predict actions from video
// @Description Upload a video URL and get predicted actions in Uzbek language
// @Tags prediction
// @Accept json
// @Produce json
// @Param request body predictActionRequest true "Request body with video URL"
// @Success 200 {object} pretictActionResponse "Returns predicted actions in Uzbek"
// @Failure 400 {object} response.BaseResponse "Invalid request body"
// @Failure 500 {object} response.BaseResponse "Internal server error or prediction service error"
// @Router /api/predict [post]
func (h *Handler) predictAction(c *gin.Context) {
	var body predictActionRequest

	err := c.ShouldBindJSON(&body)
	if err != nil {
		fmt.Println("Error binding JSON:", err)
		response.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	resp, err := h.sendRequestPredict(body.VideoURL, body.PredictApiURL)
	if err != nil {
		fmt.Println("Error sending predict request:", err)
		response.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	labels := []string{}

	predictions := resp.Predictions
	if len(predictions) == 0 {
		response.ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("no predictions found"))
		return
	}

	for i := range predictions {
		predictions[i].Label = models.ActionMapping[predictions[i].Label]
		labels = append(labels, predictions[i].Label)
	}

	go func() {
		h.usecase.Action.CreateAction(context.Background(), models.Action{
			VideoURL:         body.VideoURL,
			PredictedActions: labels,
		})
	}()

	// Set response headers
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")

	// Return success response
	c.JSON(http.StatusOK, pretictActionResponse{
		Predictions: predictions,
		VideoURL:    resp.TrackedVideoUploadURL,
	})
}

// func getActionsUzb(actions []string) []string {
// 	actionsMapping := models.ActionMapping
// 	result := []string{}

// 	for i := range actions {
// 		result = append(result, actionsMapping[actions[i]])
// 	}

// 	return result
// }

type PredictResponse struct {
	Predictions              []models.Prediction `json:"action_predictions"`
	TrackedVideoUploadStatus string              `json:"tracked_video_upload_status"`
	TrackedVideoUploadURL    string              `json:"uploaded_tracked_video_url"`
}

func (h *Handler) sendRequestPredict(link string, predictApiURL string) (PredictResponse, error) {
	// Create request body
	requestBody := map[string]string{
		"video_url":  link,
		"upload_url": "https://api.multicom.uz/api/v1/files",
	}

	// Create HTTP client
	client := &http.Client{}

	// Marshal the request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return PredictResponse{}, err
	}

	// Create the request
	req, err := http.NewRequest("POST", predictApiURL+"/predict", bytes.NewBuffer(jsonBody))
	if err != nil {
		return PredictResponse{}, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return PredictResponse{}, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return PredictResponse{}, fmt.Errorf("prediction API returned status code: %d", resp.StatusCode)
	}

	// Parse response
	var response PredictResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return PredictResponse{}, err
	}

	return response, nil
}
