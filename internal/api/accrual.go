package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/fickleDude/gophemart/internal/config"
	"github.com/fickleDude/gophemart/internal/model"
)

const MAX_RETRY = 3

type ErrorClassification int

const (
	NonRetriable ErrorClassification = iota
	Retriable
)

func Classify(err error, resp *http.Response) ErrorClassification {
	if err != nil {
		var netErr *net.OpError
		if errors.As(err, &netErr) {
			if netErr.Temporary() || netErr.Timeout() {
				return Retriable
			}
			return NonRetriable
		}
	}
	if resp.StatusCode == http.StatusBadGateway ||
		resp.StatusCode == http.StatusServiceUnavailable ||
		resp.StatusCode == http.StatusGatewayTimeout {
		return Retriable
	}
	return NonRetriable
}

func DoWithRetry(client *http.Client, req *http.Request, maxRetries int) (*http.Response, error) {
	var resp *http.Response
	var err error

	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body.Close()
	}

	backoff := 1 * time.Second

	for i := 0; i <= maxRetries; i++ {
		if bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		resp, err = client.Do(req)

		if Classify(err, resp) == NonRetriable {
			return resp, err
		}
		if resp != nil {
			resp.Body.Close()
		}

		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		case <-time.After(backoff):
		}

		//exponential backoff
		backoff *= 2
	}

	if err != nil {
		return nil, fmt.Errorf("failed after %d retries: %w", maxRetries, err)
	}
	return resp, err
}

func GetOrderAccrual(number string, client http.Client) (*model.Order, error) {
	cfg := config.GetConfig()
	url := fmt.Sprintf("%s/api/orders/%s", cfg.AccrualSystenAddress(), number)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	response, err := DoWithRetry(&client, request, MAX_RETRY)
	if err != nil {
		return nil, err
	}
	var order model.Order
	switch response.StatusCode {
	case 200:
		if err := json.NewDecoder(response.Body).Decode(&order); err != nil {
			return nil, err
		}
		defer response.Body.Close()
		return &order, nil
	case 204:
		order.Status = "NEW"
	}
	return &order, nil
}
