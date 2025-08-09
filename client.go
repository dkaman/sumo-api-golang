package sumo-api-golang

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

const (
	version          = "v0.0.1"
	defaultBaseURL   = "https://sumo-api.com"
	defaultUserAgent = "sumo-api-golang" + "/" + version
)

type Client struct {
	client    *http.Client
	BaseURL   *url.URL
	UserAgent string

	common service

	Basho    *BashoService
	Kimarite *KimariteService
	Rikishi  *RikishiService
}

type service struct {
	client *Client
}

type apiBulkResponse struct {
	Limit   int             `json:"limit,omitempty"`
	Skip    int             `json:"skip,omitempty"`
	Total   int             `json:"total,omitempty"`
	Records json.RawMessage `json:"records,omitempty"`
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	c := &Client{client: httpClient}

	c.initialize()

	return c
}

func (c *Client) initialize() {
	if c.BaseURL == nil {
		c.BaseURL, _ = url.Parse(defaultBaseURL)
	}

	if c.UserAgent == "" {
		c.UserAgent = defaultUserAgent
	}

	c.common.client = c

	c.Basho = (*BashoService)(&c.common)
	c.Kimarite = (*KimariteService)(&c.common)
	c.Rikishi = (*RikishiService)(&c.common)
}


func (c *Client) NewRequest(method string, urlStr string, body any) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, v any) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch v := v.(type) {
	case nil:

	case io.Writer:
		_, err = io.Copy(v, resp.Body)

	default:
		decErr := json.NewDecoder(resp.Body).Decode(v)
		if decErr == io.EOF {
			decErr = nil // ignore EOF errors caused by empty response body
		}

		if decErr != nil {
			err = decErr
		}
	}

	return resp, err
}
