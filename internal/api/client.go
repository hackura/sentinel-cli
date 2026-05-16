package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hackura/sentinel-cli/internal/config"
	"github.com/hackura/sentinel-cli/internal/models"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) do(method, path string, body interface{}) (*http.Response, error) {
	token, err := config.LoadToken()
	if err != nil {
		return nil, fmt.Errorf("authentication required: run 'sentinel login'")
	}

	var bodyReader io.Reader
	var bodyBytes []byte
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	maxRetries := 3
	backoff := 1 * time.Second

	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if bodyBytes != nil {
			bodyReader = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequest(method, c.BaseURL+path, bodyReader)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.HTTPClient.Do(req)
		if err == nil {
			if resp.StatusCode == http.StatusUnauthorized {
				return nil, fmt.Errorf("session expired: run 'sentinel login'")
			}
			if resp.StatusCode < 500 {
				return resp, nil
			}
			lastErr = fmt.Errorf("server error: %d", resp.StatusCode)
			resp.Body.Close()
		} else {
			lastErr = err
		}

		time.Sleep(backoff)
		backoff *= 2
	}

	return nil, fmt.Errorf("request failed after %d retries: %v", maxRetries, lastErr)
}


func (c *Client) Post(path string, body interface{}, v interface{}) error {
	resp, err := c.do("POST", path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return c.handleError(resp)
	}

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}
	return nil
}

func (c *Client) Get(path string, v interface{}) error {
	resp, err := c.do("GET", path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return c.handleError(resp)
	}

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}
	return nil
}

func (c *Client) GetUserStatus() (*models.UserStatus, error) {
	var res models.UserStatusResponse
	err := c.Get("/cli/auth/me", &res)
	if err != nil {
		return nil, err
	}
	return &res.Data, nil
}

func (c *Client) handleError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
}
