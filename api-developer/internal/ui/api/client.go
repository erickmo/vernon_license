//go:build wasm

// Package api menyediakan HTTP client untuk Vernon App internal API.
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrUnauthorized dikembalikan saat API merespons 401.
var ErrUnauthorized = fmt.Errorf("unauthorized")

// ErrorResponse adalah format JSON error dari API.
type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// Client adalah HTTP client dengan JWT auth.
type Client struct {
	baseURL string
	token   string
}

// NewClient membuat client baru dengan baseURL dan token.
func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
	}
}

// Get melakukan GET request dengan auth header.
// Decode response JSON ke result.
func (c *Client) Get(ctx context.Context, path string, result any) error {
	return c.doRequest(ctx, http.MethodGet, path, nil, result)
}

// Post melakukan POST request dengan auth header.
// Encode body sebagai JSON, decode response ke result.
func (c *Client) Post(ctx context.Context, path string, body any, result any) error {
	return c.doRequest(ctx, http.MethodPost, path, body, result)
}

// Put melakukan PUT request dengan auth header.
// Encode body sebagai JSON, decode response ke result.
func (c *Client) Put(ctx context.Context, path string, body any, result any) error {
	return c.doRequest(ctx, http.MethodPut, path, body, result)
}

// Delete melakukan DELETE request dengan auth header.
func (c *Client) Delete(ctx context.Context, path string) error {
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

// doRequest adalah helper untuk semua HTTP methods.
// Menambahkan Authorization: Bearer header dan mem-parse JSON response.
func (c *Client) doRequest(ctx context.Context, method, path string, body any, result any) error {
	var reqBody *bytes.Buffer
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("api.Client: marshal body: %w", err)
		}
		reqBody = bytes.NewBuffer(data)
	} else {
		reqBody = &bytes.Buffer{}
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("api.Client: new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("api.Client: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		// 401 selalu return ErrUnauthorized tanpa cek body
		if resp.StatusCode == 401 {
			return ErrUnauthorized
		}
		var errResp ErrorResponse
		if decErr := json.NewDecoder(resp.Body).Decode(&errResp); decErr == nil && errResp.Error.Message != "" {
			return fmt.Errorf("api error %d: %s — %s", resp.StatusCode, errResp.Error.Code, errResp.Error.Message)
		}
		if resp.StatusCode == 403 {
			return fmt.Errorf("Anda tidak memiliki izin untuk mengakses resource ini")
		}
		return fmt.Errorf("api error: status %d", resp.StatusCode)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("api.Client: decode response: %w", err)
		}
	}

	return nil
}
