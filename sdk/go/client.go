// Package freepastebin provides a dependency-free client for freepastebin.com.
package freepastebin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}
type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func (r *Response) Text() string     { return string(r.Body) }
func (r *Response) JSON(v any) error { return json.Unmarshal(r.Body, v) }

type APIError struct {
	StatusCode int
	Body       []byte
}

func (e *APIError) Error() string {
	body := e.Body
	if len(body) > 500 {
		body = body[:500]
	}
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, body)
}

type Options struct {
	Query          map[string]any
	Method         string
	Headers        map[string]string
	Body           any
	RawBody        []byte
	AllowHTTPError bool
}

func New() *Client {
	return &Client{BaseURL: "https://freepastebin.com", HTTPClient: &http.Client{Timeout: 45 * time.Second}}
}
func (c *Client) URL(path string, query map[string]any) (string, error) {
	if !strings.HasPrefix(path, "/") || strings.HasPrefix(path, "//") {
		return "", fmt.Errorf("use an API path beginning with one slash")
	}
	base := c.BaseURL
	if base == "" {
		base = "https://freepastebin.com"
	}
	u, err := url.Parse(strings.TrimRight(base, "/") + path)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for k, v := range query {
		if v == nil {
			continue
		}
		switch x := v.(type) {
		case string:
			q.Set(k, x)
		default:
			b, e := json.Marshal(v)
			if e != nil {
				return "", e
			}
			q.Set(k, string(b))
		}
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// Request buffers the response body and leaves redirects visible to callers.
func (c *Client) Request(ctx context.Context, path string, o Options) (*Response, error) {
	target, err := c.URL(path, o.Query)
	if err != nil {
		return nil, err
	}
	method := o.Method
	if method == "" {
		method = "GET"
	}
	body := o.RawBody
	if o.Body != nil {
		body, err = json.Marshal(o.Body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, target, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	for k, v := range o.Headers {
		req.Header.Set(k, v)
	}
	if o.Body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 45 * time.Second}
	}
	copyClient := *client
	copyClient.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	res, err := copyClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	out := &Response{StatusCode: res.StatusCode, Header: res.Header, Body: data}
	if !o.AllowHTTPError && (res.StatusCode < 200 || res.StatusCode >= 300) {
		return out, &APIError{StatusCode: res.StatusCode, Body: data}
	}
	return out, nil
}
func copyMap(in map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range in {
		out[k] = v
	}
	return out
}

// Create defaults to POST; pass method GET for small, non-sensitive pastes.
func (c *Client) Create(ctx context.Context, content string, params map[string]any, method string) (*Response, error) {
	q := copyMap(params)
	q["content"] = content
	if method == "GET" {
		return c.Request(ctx, "/api/pastes", Options{Query: q})
	}
	if method == "" {
		method = "POST"
	}
	return c.Request(ctx, "/api/pastes", Options{Method: method, Body: q})
}
func (c *Client) Get(ctx context.Context, id, readToken string) (*Response, error) {
	h := map[string]string{}
	if readToken != "" {
		h["Authorization"] = "Bearer " + readToken
	}
	return c.Request(ctx, "/api/pastes/"+url.PathEscape(id), Options{Headers: h})
}
func (c *Client) Raw(ctx context.Context, id, readToken string) (*Response, error) {
	h := map[string]string{}
	if readToken != "" {
		h["Authorization"] = "Bearer " + readToken
	}
	return c.Request(ctx, "/raw/"+url.PathEscape(id), Options{Headers: h})
}
func (c *Client) Delete(ctx context.Context, id, token string) error {
	_, err := c.Request(ctx, "/api/pastes/"+url.PathEscape(id)+"/delete", Options{Method: "POST", Body: map[string]any{"delete_token": token}})
	return err
}
