package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/xaionaro-go/homenet-server/helpers"
)

func returnSuccess(ctx *gin.Context, data interface{}) {
	helpers.ReturnSuccess(ctx, data)
}

func returnError(ctx *gin.Context, err error) {
	helpers.ReturnError(ctx, err)
	ctx.Abort()
}
