package yourls

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/devusSs/minls/internal/log"
)

func (c *Client) Shorten(ctx context.Context, input string) (string, error) {
	if ctx == nil {
		return "", errors.New("nil context")
	}

	uid, err := uuid.NewUUID()
	if err != nil {
		return "", fmt.Errorf("could not create uuid for keyword: %w", err)
	}

	log.Debug(
		"yourls - *client.Shorten",
		slog.String("action", "uuid_newuuid"),
		slog.Any("uid", uid),
	)

	v := make(map[string]string)
	v["signature"] = c.signature
	v["action"] = "shorturl"
	v["format"] = "json"
	v["url"] = input
	v["title"] = defaultUploadTitle
	v["keyword"] = uid.String()

	log.Debug("yourls - *client.Shorten", slog.String("action", "set_values"), slog.Any("v", v))

	resp, err := c.do(ctx, v)
	if err != nil {
		return "", fmt.Errorf("client.do(): %w", err)
	}
	defer resp.Body.Close()

	log.Debug("yourls - *client.Shorten", slog.String("action", "got_resp"), slog.Any("resp", resp))

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unwanted status code: %d (%s)", resp.StatusCode, resp.Status)
	}

	res := &shortenURLResponse{}
	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return "", fmt.Errorf("could not decode response: %w", err)
	}

	log.Debug(
		"yourls - *client.Shorten",
		slog.String("action", "decoded_resp"),
		slog.Any("res", res),
	)

	return res.Shorturl, nil
}

const defaultUploadTitle = "Uploaded using minls by devusSs"

type shortenURLResponse struct {
	URL struct {
		Keyword string `json:"keyword"`
		URL     string `json:"url"`
		Title   string `json:"title"`
		Date    string `json:"date"`
		IP      string `json:"ip"`
	} `json:"url"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Title      string `json:"title"`
	Shorturl   string `json:"shorturl"`
	StatusCode string `json:"statusCode"`
}
