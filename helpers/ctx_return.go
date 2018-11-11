package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/xaionaro-go/homenet-server/errors"
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
	case errors.BadRequest:
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":            "error",
			"error_type":        "bad_request",
			"error_description": err.Error(),
		})
	default:
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":     "error",
			"error_type": "internal_server_error",
		})
	}
}
