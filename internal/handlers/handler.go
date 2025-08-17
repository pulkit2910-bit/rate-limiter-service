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
	ConfigHandler(ctx *gin.Context)
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
		fmt.Printf("Error running Lua script: %v\n", err)
        return
    }

	fmt.Printf("Lua script result: %v\n", result)

    ctx.JSON(http.StatusOK, gin.H{"result": result})
}

func (h *Handler) ConfigHandler(ctx *gin.Context) {
	var req struct {
        UserId     string `json:"userId"`
        Capacity   string `json:"capacity"`
        RefillRate string `json:"refillRate"`
    }

    if err := ctx.BindJSON(&req); err != nil {
		fmt.Printf("Error binding request body: %v\n", err)
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }

    if req.UserId == "" || req.Capacity == "" || req.RefillRate == "" {
		fmt.Printf("Invalid request: userId, capacity, and refillRate are required\n")
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "userId, capacity, and refillRate are required"})
        return
    }

    if err := h.service.RateLimitingConfig(ctx, req.UserId, req.Capacity, req.RefillRate); err != nil {
		fmt.Printf("Error setting rate limiting config: %v\n", err)
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set rate limiting config"})
        return
    }

	fmt.Printf("Rate limiting config set for user %s: capacity=%s, refillRate=%s\n", req.UserId, req.Capacity, req.RefillRate)
	ctx.JSON(http.StatusOK, gin.H{"message": "Rate limiting config set successfully"})
}