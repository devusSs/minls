package yourls

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	endpoint  string
	signature string
	client    *http.Client
}

func NewClient(endpoint string, signature string) *Client {
	if !strings.Contains(endpoint, "https://") && !strings.Contains(endpoint, "http://") {
		endpoint = "http://" + endpoint
	}

	if !strings.Contains(endpoint, "https://") {
	}

	if !strings.Contains(endpoint, "/yourls-api.php") {
	}

	return &Client{
		endpoint:  endpoint,
		signature: signature,
		client:    http.DefaultClient,
	}
}

func (c *Client) do(
	ctx context.Context,
	values map[string]string,
) (*http.Response, error) {
	if ctx == nil {
		return nil, errors.New("nil context")
	}

	if values == nil {
		return nil, errors.New("nil values")
	}

	u, err := url.Parse(c.endpoint)
	if err != nil {
		return nil, fmt.Errorf("could not parse endpoint: %w", err)
	}

	_, ok := values["signature"]
	if !ok {
		values["signature"] = c.signature
	}

	vl := url.Values{}
	for k, v := range values {
		vl.Set(k, v)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		u.String(),
		strings.NewReader(vl.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("could not build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not get response: %w", err)
	}

	return resp, nil
}
