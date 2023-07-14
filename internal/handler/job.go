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

type JobHandler interface {
	Launch(c *gin.Context)
}

type jobHandler struct {
	s service.JobService
}

func NewJobHandler(s service.JobService) JobHandler {
	return &jobHandler{s: s}
}

// Launch godoc
// @Summary Launches a job
// @Description Launches a job
// @Tags jobs
// @Accept json
// @Produce json
// @Param workerID path string true "Worker ID"
// @Param input body map[string]any true "Input"
// @Success 200 {object} LaunchJobResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /workers/{workerID}/jobs [post]
func (h *jobHandler) Launch(c *gin.Context) {
	workerID := c.Param("workerID")
	var input map[string]any
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("failed to bind input. error: %v", err)
		c.JSON(http.StatusBadRequest, toErrorResponse(exception.JobError{
			Code:    exception.InvalidRequest,
			Message: fmt.Sprintf("invalid request. error: %v", err),
			Err:     err,
		}))
		return
	}

	jobID, err := h.s.LaunchJob(c, workerID, input)
	if err != nil {
		log.Printf("failed to launch job. error: %v", err)
		c.JSON(http.StatusInternalServerError, toErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.LaunchJobResponse{
		JobID: jobID,
	})
}
