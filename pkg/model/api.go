package model

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Worker struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Topic string `json:"topic"`
}

type CreateWorkerRequest struct {
	Name  string `json:"name"`
	Topic string `json:"topic"`
}

type CreateWorkerResponse struct {
	ID string `json:"id"`
}
