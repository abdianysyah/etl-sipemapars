package callback

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	BaseURL		string
	Secret		string
	Client		*http.Client
}

func New(baseURL, secret string) *Client {
	return &Client {
		BaseURL: strings.TrimRight(baseURL, "/"),
		Secret:	secret,
		Client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) post(path string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err		
	}

	req, err := http.NewRequest(http.MethodPost, c.BaseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header .Set("Content-Type", "application/json")
	if c.Secret != "" {
		req.Header.Set("X-ETL-SECRET", c.Secret)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("callback %s failed: %s", path, resp.Status)
	}

	return nil
}

func (c *Client) Progress(jobUUID string, progress int, currentStep string) error {
	return c.post("/api/etl/progress", map[string]any{
		"job_uuid": jobUUID,
		"progress": progress,
		"current_step": currentStep,
	})
}

func (c *Client) Log(jobUUID string, level int, step, message string) error {
	return c.post("/api/etl/log", map[string]any{
		"job_uuid": jobUUID,
		"level": level,
		"step": step,
		"message": message,
	})
}

func (c *Client) Finish(jobUUID string, message string) error {
	return c.post("/api/etl/finish", map[string]any{
		"job_uuid": jobUUID,
		"message": message,
	})
}

func (c *Client) Failed(jobUUID string, step string, message string) error {
	return c.post("/api/etl/failed", map[string]any{
		"job_uuid": jobUUID,
		"step": step,
		"message": message,
	})
}