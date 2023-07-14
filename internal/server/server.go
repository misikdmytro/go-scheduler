package server

import (
	"github.com/gin-gonic/gin"
	"github.com/misikdmytro/go-job-runner/docs"
	"github.com/misikdmytro/go-job-runner/internal/handler"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewEngine(
	wh handler.WorkerHandler,
	jh handler.JobHandler,
	hh handler.HealthHandler,
) *gin.Engine {
	e := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"

	w := e.Group("/workers")
	{
		w.GET("/:id", wh.Get)
		w.PUT("/", wh.Create)
		w.DELETE("/:id", wh.Delete)

		w.POST("/:id/jobs", jh.Launch)
	}

	h := e.Group("/health")
	{
		h.GET("/", hh.Health)
	}

	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return e
}
