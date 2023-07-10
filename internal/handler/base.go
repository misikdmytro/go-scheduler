package handler

import (
	"errors"

	"github.com/misikdmytro/go-job-runner/internal/exception"
	"github.com/misikdmytro/go-job-runner/pkg/model"
)

func toErrorResponse(err error) any {
	var e exception.JobError
	if ok := errors.As(err, &e); ok {
		return model.ErrorResponse{
			Code:    int(e.Code),
			Message: e.Message,
		}
	}

	return model.ErrorResponse{
		Code:    int(exception.UnknownError),
		Message: "unknown error",
	}
}
