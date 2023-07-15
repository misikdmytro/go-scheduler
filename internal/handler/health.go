package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/misikdmytro/go-job-runner/internal/service"
)

type HealthHandler interface {
	Health(c *gin.Context)
}

type healthHandler struct {
	s service.HealthService
}

func NewHealthHandler(s service.HealthService) HealthHandler {
	return &healthHandler{s: s}
}

// Health godoc
// @Summary Health check
// @Description Health check
// @Tags health
// @Accept json
// @Produce json
// @Success 204
// @Failure 500
// @Router /health [get]
func (h *healthHandler) Health(c *gin.Context) {
	if err := h.s.Health(c); err != nil {
		log.Printf("failed to check health. error: %v", err)
		c.JSON(http.StatusInternalServerError, toErrorResponse(err))
		return
	}

	c.Status(http.StatusNoContent)
}
