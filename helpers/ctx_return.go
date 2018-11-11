package helpers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Sirupsen/logrus"
)

func ReturnSuccess(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"status": "OK",
		"result": data,
	})
	return
}

func ReturnError(ctx *gin.Context, err error) {
	switch err.(type) {
	case interface{IsNotFound()}:
		ctx.JSON(http.StatusNotFound, map[string]interface{}{
			"status":            "error",
			"error_type":        "not_found",
			"error_description": err.Error(),
		})
	case interface{IsBadRequest()}:
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":            "error",
			"error_type":        "bad_request",
			"error_sub_type":    fmt.Sprintf("%T", err),
			"error_description": err.Error(),
		})
	default:
		logrus.Errorf("Got error: %T: %v", err, err)
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":     "error",
			"error_type": "internal_server_error",
		})
	}
}
