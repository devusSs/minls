package yourls

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/devusSs/minls/internal/log"
)

type Client struct {
	endpoint  string
	signature string
	client    *http.Client
}

func NewClient(endpoint string, signature string) *Client {
	if !strings.Contains(endpoint, "https://") && !strings.Contains(endpoint, "http://") {
		log.Warn(
			"yourls - NewClient",
			slog.String("warn", "endpoint does not contain schema, adding 'http://'"),
		)
		endpoint = "http://" + endpoint
	}

	if !strings.Contains(endpoint, "https://") {
		log.Warn("yourls - NewClient", slog.String("warn", "endpoint is not secure"))
	}

	if !strings.Contains(endpoint, "/yourls-api.php") {
		log.Warn(
			"yourls - NewClient",
			slog.String("warn", "endpoint does not contains typical api path"),
		)
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

	log.Debug("yourls - *client.do", slog.String("action", "url_parse"), slog.Any("parsed_url", u))

	_, ok := values["signature"]
	if !ok {
		log.Warn(
			"yourls - *client.do",
			slog.String("action", "check_values"),
			slog.String("warn", "signature not in values, adding from client"),
		)
		values["signature"] = c.signature
	}

	vl := url.Values{}
	for k, v := range values {
		vl.Set(k, v)
	}

	log.Debug("yourls - *client.do", slog.String("action", "set_values"), slog.Any("vl", vl))

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

	log.Debug("yourls - *client.do", slog.String("action", "got_response"), slog.Any("resp", resp))

	return resp, nil
}
