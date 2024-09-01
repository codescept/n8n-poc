package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"example.com/n8n-poc/config"
	"example.com/n8n-poc/errors"
)

type Client interface {
	Do(ctx context.Context, method string, path string, payload []byte) (*http.Response, *errors.ErrResponse)
}

type impl struct{}

const (
	resourceURL = "%s/%s"
)

func New() Client {
	return &impl{}
}

func (i impl) Do(ctx context.Context, method string, path string, payload []byte) (*http.Response, *errors.ErrResponse) {
	conf := config.Get()
	url := fmt.Sprintf(resourceURL, conf.N8N_BASE_URL, path)

	httpReq, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, &errors.ErrResponse{
			Message: err.Error(),
		}
	}

	httpReq.Header.Add("X-N8N-API-KEY", conf.N8N_API_KEY)
	httpReq.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, &errors.ErrResponse{
			Message: err.Error(),
		}
	}

	if resp.StatusCode >= http.StatusBadRequest {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, &errors.ErrResponse{
				Message: err.Error(),
				Status:  resp.StatusCode,
			}
		}
		var errResp errors.ErrResponse
		if err = json.Unmarshal(body, &errResp); err != nil {
			return nil, &errors.ErrResponse{
				Message: err.Error(),
				Status:  resp.StatusCode,
			}
		}

		errResp.Status = resp.StatusCode

		return nil, &errResp
	}

	return resp, nil
}
