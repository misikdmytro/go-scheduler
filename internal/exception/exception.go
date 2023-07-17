package exception

import "fmt"

type JobErrorCode int

const (
	InvalidRequest JobErrorCode = iota + 40000
	NothingInserted
	WorkerNotFound
	WorkerAlreadyExists

	UnknownError JobErrorCode = 50000
	FailedToAppendJobStatus
	UnhealthService
)

type JobError struct {
	Code    JobErrorCode
	Message string
	Err     error
}

func (e JobError) Error() string {
	return fmt.Sprintf(
		"job error. code: %d. message: %s. internal error: %v",
		e.Code,
		e.Message,
		e.Err,
	)
}
