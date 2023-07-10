package server

import (
	"github.com/gin-gonic/gin"
	"github.com/misikdmytro/go-job-runner/internal/handler"
)

func NewEngine(
	wh handler.WorkerHandler,
	hh handler.HealthHandler,
) *gin.Engine {
	e := gin.Default()

	w := e.Group("/workers")
	{
		w.GET("/:id", wh.Get)
		w.PUT("/", wh.Create)
		w.DELETE("/:id", wh.Delete)
	}

	h := e.Group("/health")
	{
		h.GET("/", hh.Health)
	}

	return e
}
