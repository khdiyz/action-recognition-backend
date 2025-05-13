package handler

import (
	"action-detector-backend/models"
	"action-detector-backend/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type getActionsResponse struct {
	Actions []models.Action `json:"actions"`
}

// @Summary Get all actions
// @Description Get a list of all previously predicted actions
// @Tags actions
// @Produce json
// @Success 200 {object} getActionsResponse "Returns the list of actions"
// @Failure 500 {object} response.BaseResponse "Internal server error"
// @Router /api/actions [get]
func (h *Handler) getActions(c *gin.Context) {
	actions, err := h.usecase.GetActions(c.Request.Context())
	if err != nil {
		response.ErrorResponse(c, 500, err)
		return
	}

	c.JSON(http.StatusOK, getActionsResponse{
		Actions: actions,
	})
}

// @Summary Delete all actions
// @Description Delete all previously predicted actions
// @Tags actions
// @Produce json
// @Success 200 {object} response.BaseResponse "Successfully deleted all actions"
// @Failure 500 {object} response.BaseResponse "Internal server error"
// @Router /api/actions [delete]
func (h *Handler) deleteActions(c *gin.Context) {
	err := h.usecase.DeleteActions(c.Request.Context())
	if err != nil {
		response.ErrorResponse(c, 500, err)
		return
	}

	response.SuccessResponse(c, "All actions successfully deleted")
}
