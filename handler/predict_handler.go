package handler

import (
	"action-detector-backend/models"
	"action-detector-backend/pkg/response"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// pretictActionResponse represents the response structure for action prediction
// @Description Response containing predicted actions in Uzbek language
type pretictActionResponse struct {
	// @Description List of predicted actions in Uzbek language
	Predictions []string `json:"predictions"`
}

// predictActionRequest represents the request body for predicting actions
// @Description Request containing the URL of the video to analyze
type predictActionRequest struct {
	// @Description URL of the video to analyze
	VideoURL string `json:"video_url" binding:"required"`
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
		response.ErrorResponse(c, 400, err)
		return
	}

	predictions, err := h.sendRequestPredict(body.VideoURL)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	result := getActionsUzb(predictions)

	fmt.Println(result)

	// Return success response
	c.JSON(http.StatusOK, pretictActionResponse{
		Predictions: result,
	})
}

func getActionsUzb(actions []string) []string {
	actionsMapping := models.ActionMapping
	result := []string{}

	for i := range actions {
		result = append(result, actionsMapping[actions[i]])
	}

	return result
}

func (h *Handler) sendRequestPredict(link string) ([]string, error) {
	// Create request body
	requestBody := map[string]string{
		"video_url": link,
	}

	// Create HTTP client
	client := &http.Client{}

	// Marshal the request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	// Create the request
	req, err := http.NewRequest("POST", h.cfg.PredictApiURL+"/predict", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("prediction API returned status code: %d", resp.StatusCode)
	}

	// Parse response
	var response struct {
		Predictions []string `json:"predictions"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Predictions, nil
}
