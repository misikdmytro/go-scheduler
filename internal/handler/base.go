package handler

import (
	"errors"
	"net/http"

	"github.com/misikdmytro/go-job-runner/internal/exception"
	"github.com/misikdmytro/go-job-runner/pkg/model"
)

func toErrorResponse(err error) (int, any) {
	var e exception.JobError
	if ok := errors.As(err, &e); ok {
		statusCode := int(e.Code) / 100
		return statusCode, model.ErrorResponse{
			Code:    int(e.Code),
			Message: e.Message,
		}
	}

	return http.StatusInternalServerError, model.ErrorResponse{
		Code:    int(exception.UnknownError),
		Message: "unknown error",
	}
}
