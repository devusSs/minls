package cli

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/devusSs/minls/internal/log"
	"github.com/devusSs/minls/internal/storage"
)

func List() error {
	err := initialize()
	if err != nil {
		return fmt.Errorf("could not initialize: %w", err)
	}

	log.Debug("cli - List", slog.String("action", "initialize"))

	data, err := storage.ReadData()
	if err != nil {
		return fmt.Errorf("could not read data from storage: %w", err)
	}

	log.Debug("cli - List", slog.String("action", "storage_read_data"), slog.Any("data", data))

	if len(data.Entries) == 0 {
		fmt.Println("No data to be displayed yet.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "ID\tTimestamp\tMinio ID\tYOURLS ID")

	for _, entry := range data.Entries {
		var minioURL *url.URL
		var yourlsURL *url.URL

		minioURL, err = url.Parse(entry.MinioLink)
		if err != nil {
			return fmt.Errorf("malformed minio url: %w", err)
		}

		yourlsURL, err = url.Parse(entry.YOURLSLink)
		if err != nil {
			return fmt.Errorf("mailformed yourls url: %w", err)
		}

		fmt.Fprintf(
			w,
			"%d\t%s\t%s\t%s\n",
			entry.ID,
			entry.Timestamp.Format("2006-01-02 15:04:05"),
			filepath.Base(minioURL.Path),
			filepath.Base(yourlsURL.Path),
		)
	}

	return w.Flush()
}
