package handler

import (
	"fmt"
	"net/http"
	"os"
	"time"

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
    userId := ctx.Query("userId")
	apiPath := ctx.Query("api")
    if userId == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
        return
    }
    if apiPath == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "api is required"})
        return
    }
	
    scriptBytes, err := os.ReadFile("scripts/script.lua")
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read Lua script"})
        return
    }
    
	currentTimeStamp := time.Now().UnixMilli()
    script := string(scriptBytes)
    keys := []string{userId, apiPath}
	args := []string{"current_tokens", "last_updated", "capacity", "refill_rate", fmt.Sprintf("%d", currentTimeStamp)}

    result, err := h.service.RunLuaScript(ctx, script, keys, args)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"result": result})
}