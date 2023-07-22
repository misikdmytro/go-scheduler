package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
	"github.com/misikdmytro/go-job-runner/internal/exception"
	"github.com/misikdmytro/go-job-runner/internal/model"
	"github.com/misikdmytro/go-job-runner/internal/service"
	pkgmodel "github.com/misikdmytro/go-job-runner/pkg/model"
)

type JobHandler interface {
	Launch(*gin.Context)
	JobStatuses(*gin.Context)
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
	var req pkgmodel.LaunchJobRequest
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

	c.JSON(http.StatusOK, pkgmodel.LaunchJobResponse{
		JobID: jobID,
	})
}

// JobStatuses godoc
// @Summary Get job statuses
// @Description Get job statuses
// @Tags jobs
// @Accept json
// @Produce json
// @Param jobID path string true "Job ID"
// @Success 200 {array} model.JobStatusAPI
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /jobs/{jobID}/statuses [get]
func (h *jobHandler) JobStatuses(c *gin.Context) {
	jobID := c.Param("id")

	statuses, err := h.s.GetJobStatuses(c, jobID)
	if err != nil {
		log.Printf("failed to get job statuses. error: %v", err)
		c.JSON(toErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, pkgmodel.JobStatusesResponse{
		Statuses: pie.Map(statuses, func(s model.JobStatus) pkgmodel.JobStatusAPI {
			return pkgmodel.JobStatusAPI{
				Message:   s.Message,
				Timestamp: s.Timestamp,
				Output:    s.Output,
			}
		}),
	})
}
