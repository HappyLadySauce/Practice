package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// 统一错误响应函数
func respondWithError(ctx *gin.Context, statusCode int, message string) {
	ctx.String(statusCode, message)
}

// 统一成功响应函数
func respondWithSuccess(ctx *gin.Context, message string) {
	ctx.String(http.StatusOK, message)
}