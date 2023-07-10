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
