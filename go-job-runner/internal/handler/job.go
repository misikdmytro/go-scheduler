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
// @Summary Launch job
// @Description Launch job
// @Tags jobs
// @Accept json
// @Produce json
// @Param request body model.LaunchJobRequest true "Launch job request"
// @Success 200 {object} model.LaunchJobResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /jobs [post]
func (h *jobHandler) Launch(c *gin.Context) {
	var req model.LaunchJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("failed to bind request. error: %v", err)
		c.JSON(toErrorResponse(exception.JobError{
			Code:    exception.InvalidRequest,
			Message: fmt.Sprintf("invalid request. error: %v", err),
			Err:     err,
		}))
		return
	}

	jobID, err := h.s.LaunchJob(c, req.WorkerID, req.Input)
	if err != nil {
		log.Printf("failed to launch job. error: %v", err)
		c.JSON(toErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.LaunchJobResponse{
		JobID: jobID,
	})
}
