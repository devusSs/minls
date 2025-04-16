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

	log.Info("yourls - shorten", slog.String("input", input), slog.String("uuid", uid.String()))

	v := make(map[string]string)
	v["signature"] = c.signature
	v["action"] = "shorturl"
	v["format"] = "json"
	v["url"] = input
	v["title"] = defaultUploadTitle
	v["keyword"] = uid.String()

	log.Debug("yourls - shorten", slog.Any("values", v))

	resp, err := c.do(ctx, v)
	if err != nil {
		return "", fmt.Errorf("client.do(): %w", err)
	}
	defer resp.Body.Close()

	log.Debug(
		"yourls - shorten",
		slog.Int("resp_code", resp.StatusCode),
		slog.String("resp_status", resp.Status),
	)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unwanted status code: %d (%s)", resp.StatusCode, resp.Status)
	}

	res := &shortenURLResponse{}
	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return "", fmt.Errorf("could not decode response: %w", err)
	}

	log.Debug("yourls - shorten", slog.Any("res", res))
	log.Info("yourls - shorten", slog.String("shortened", res.Shorturl))

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
