package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/misikdmytro/go-job-runner/internal/exception"
	"github.com/misikdmytro/go-job-runner/internal/service"
	"github.com/misikdmytro/go-job-runner/pkg/model"
)

type WorkerHandler interface {
	Create(*gin.Context)
}

type workerHandler struct {
	s service.WorkerService
}

func NewWorkerHandler(s service.WorkerService) WorkerHandler {
	return &workerHandler{s: s}
}

func (h *workerHandler) Create(c *gin.Context) {
	var req model.CreateWorkerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("failed to bind request. error: %v", err)
		c.JSON(http.StatusBadRequest, toErrorResponse(exception.JobError{
			Code:    exception.InvalidRequest,
			Message: fmt.Sprintf("invalid request. error: %v", err),
			Err:     err,
		}))
		return
	}

	id, err := h.s.Create(c, req.Name, req.Description)
	if err != nil {
		log.Printf("failed to create worker. error: %v", err)
		c.JSON(http.StatusInternalServerError, toErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, model.CreateWorkerResponse{ID: id})
}
