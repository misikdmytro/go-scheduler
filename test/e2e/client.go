package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/misikdmytro/go-job-runner/pkg/model"
)

type client struct {
	baseAddress string
}

func newClient() *client {
	return &client{baseAddress: "http://localhost:4001"}
}

func (c *client) CreateWorker(name string, desc string) (model.CreateWorkerResponse, error) {
	cl := http.Client{}

	cwr := model.CreateWorkerRequest{
		Name:        name,
		Description: desc,
	}

	b, err := json.Marshal(cwr)
	if err != nil {
		return model.CreateWorkerResponse{}, err
	}

	r := bytes.NewReader(b)
	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/workers", c.baseAddress),
		r,
	)

	if err != nil {
		return model.CreateWorkerResponse{}, err
	}

	resp, err := cl.Do(req)
	if err != nil {
		return model.CreateWorkerResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return model.CreateWorkerResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var cwresp model.CreateWorkerResponse
	if err := json.NewDecoder(resp.Body).Decode(&cwresp); err != nil {
		return model.CreateWorkerResponse{}, err
	}

	return cwresp, nil
}

func (c *client) GetWorker(id string) (model.GetWorkerResponse, error) {
	cl := http.Client{}

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/workers/%s", c.baseAddress, id),
		nil,
	)

	if err != nil {
		return model.GetWorkerResponse{}, err
	}

	resp, err := cl.Do(req)
	if err != nil {
		return model.GetWorkerResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.GetWorkerResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var gwresp model.GetWorkerResponse
	if err := json.NewDecoder(resp.Body).Decode(&gwresp); err != nil {
		return model.GetWorkerResponse{}, err
	}

	return gwresp, nil
}

func (c *client) DeleteWorker(id string) (model.DeleteWorkerResponse, error) {
	cl := http.Client{}

	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/workers/%s", c.baseAddress, id),
		nil,
	)

	if err != nil {
		return model.DeleteWorkerResponse{}, err
	}

	resp, err := cl.Do(req)
	if err != nil {
		return model.DeleteWorkerResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.DeleteWorkerResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var dwresp model.DeleteWorkerResponse
	if err := json.NewDecoder(resp.Body).Decode(&dwresp); err != nil {
		return model.DeleteWorkerResponse{}, err
	}

	return dwresp, nil
}

func (c *client) LaunchJob(workerID string, body map[string]any) (model.LaunchJobResponse, error) {
	cl := http.Client{}

	b, err := json.Marshal(body)
	if err != nil {
		return model.LaunchJobResponse{}, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/workers/%s/jobs", c.baseAddress, workerID),
		bytes.NewReader(b),
	)

	if err != nil {
		return model.LaunchJobResponse{}, err
	}

	resp, err := cl.Do(req)
	if err != nil {
		return model.LaunchJobResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.LaunchJobResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var ljresp model.LaunchJobResponse
	if err := json.NewDecoder(resp.Body).Decode(&ljresp); err != nil {
		return model.LaunchJobResponse{}, err
	}

	return ljresp, nil
}
