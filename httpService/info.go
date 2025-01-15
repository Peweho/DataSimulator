package httpService

import (
	"github.com/gin-gonic/gin"
)

func HttpMsg(ctx *gin.Context, code int, msg string) {
	ctx.JSON(code, gin.H{
		"msg": msg,
	})
}

func HttpResultMsg(ctx *gin.Context, code int, msg string) {
	ctx.JSON(code, gin.H{
		"msg": msg,
	})
}
