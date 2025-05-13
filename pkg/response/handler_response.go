package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	CreatedMessage = "created"
	SuccessMessage = "success"
)

type BaseResponse struct {
	Message string `json:"message"`
}

type IdResponse struct {
	Id string `json:"id"`
}

// Helper function to return an error response
func ErrorResponse(c *gin.Context, status int, err error) {
	c.JSON(status, BaseResponse{
		Message: err.Error(),
	})
}

func SuccessResponse(c *gin.Context, message string) {
	c.JSON(http.StatusOK, BaseResponse{
		Message: message,
	})
}

// Converts a gRPC error into an HTTP response
func FromError(c *gin.Context, serviceError error) {
	st, _ := status.FromError(serviceError)
	err := st.Message()

	switch st.Code() {
	case codes.NotFound:
		ErrorResponse(c, http.StatusNotFound, errors.New(err))
	case codes.InvalidArgument:
		ErrorResponse(c, http.StatusBadRequest, errors.New(err))
	case codes.Unavailable:
		ErrorResponse(c, http.StatusUnavailableForLegalReasons, errors.New(err))
	case codes.AlreadyExists:
		ErrorResponse(c, http.StatusBadRequest, errors.New(err))
	case codes.Unauthenticated:
		ErrorResponse(c, http.StatusUnauthorized, errors.New(err))
	default:
		ErrorResponse(c, http.StatusInternalServerError, errors.New(err))
	}
}
