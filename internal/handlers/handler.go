package handler

import (
	"github.com/gin-gonic/gin"
	service "github.com/pulkit2910-bit/rate-limiter-service/internal/services/limiter"
)

type Handler struct {
    service service.LimiterService
}

type LimiterHandler interface {
	CheckHandler(ctx *gin.Context)
}

func NewHandler(service service.LimiterService) LimiterHandler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CheckHandler(ctx *gin.Context) {
    // userId := ctx.Param("userId")


}