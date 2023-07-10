package server

import (
	"github.com/gin-gonic/gin"
	"github.com/misikdmytro/go-job-runner/internal/handler"
)

func NewEngine(
	wh handler.WorkerHandler,
	hh handler.HealthHandler,
) *gin.Engine {
	engine := gin.Default()

	w := engine.Group("/worker")
	{
		w.PUT("/", wh.Create)
	}

	h := engine.Group("/health")
	{
		h.GET("/", hh.Health)
	}

	return engine
}
