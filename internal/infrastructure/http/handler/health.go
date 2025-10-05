package handler

import (
	"net/http"

	"github.com/ambiyansyah-risyal/go-base-app/internal/usecase"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/logger"
	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	healthUsecase usecase.HealthUsecase
	logger        *logger.Logger
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(healthUsecase usecase.HealthUsecase, log *logger.Logger) *HealthHandler {
	return &HealthHandler{
		healthUsecase: healthUsecase,
		logger:        log,
	}
}

// Health godoc
// @Summary Health check
// @Description Get application health status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Health status"
// @Failure 503 {object} map[string]interface{} "Service unavailable"
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	health := h.healthUsecase.GetHealth(c.Request.Context())
	
	statusCode := http.StatusOK
	if health["status"] != "ok" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, health)
}

// Live godoc
// @Summary Liveness probe
// @Description Check if the application is alive
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Application is alive"
// @Router /health/live [get]
func (h *HealthHandler) Live(c *gin.Context) {
	if h.healthUsecase.IsLive(c.Request.Context()) {
		c.JSON(http.StatusOK, gin.H{
			"status": "alive",
		})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "dead",
		})
	}
}

// Ready godoc
// @Summary Readiness probe
// @Description Check if the application is ready to serve requests
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Application is ready"
// @Failure 503 {object} map[string]interface{} "Application is not ready"
// @Router /health/ready [get]
func (h *HealthHandler) Ready(c *gin.Context) {
	if h.healthUsecase.IsReady(c.Request.Context()) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
		})
	}
}