package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func Init() error {
	err := createStorageDirIfNotExists()
	if err != nil {
		return fmt.Errorf("could not create storage dir: %w", err)
	}

	err = createStorageFileIfNotExists()
	if err != nil {
		return fmt.Errorf("could not create storage file: %w", err)
	}

	currentData, err = readCurrentDataFromFile()
	if err != nil {
		return fmt.Errorf("could not read current data from storage file: %w", err)
	}

	// TODO: run this in goroutine for performance?
	err = cleanData()
	if err != nil {
		return fmt.Errorf("could not clean data: %w", err)
	}

	return nil
}

func WriteEntry(entry *DataEntry) error {
	err := entry.validate()
	if err != nil {
		return fmt.Errorf("could not validate entry: %w", err)
	}

	currentData.Entries = append(currentData.Entries, entry)

	err = writeCurrentData()
	if err != nil {
		return fmt.Errorf("could not write current data: %w", err)
	}

	return nil
}

func RemoveStorageDir() error {
	err := os.RemoveAll(storageDir)
	if err != nil {
		return fmt.Errorf("could not os.RemoveAll: %w", err)
	}

	return nil
}

func ReadData() (*Data, error) {
	return readCurrentDataFromFile()
}

type Data struct {
	Entries []*DataEntry `json:"entries"`
}

// DataEntry is an entry in our data.
// When creating a DataEntry it is not required to
// set an ID or timestamp, they will be set
// automatically if not provided.
type DataEntry struct {
	ID         int       `json:"id"`
	Timestamp  time.Time `json:"timestamp"`
	MinioLink  string    `json:"minio_link"`
	YOURLSLink string    `json:"yourls_link"`
}

func (e *DataEntry) validate() error {
	if e.ID == 0 {
		latestID, err := findLatestID()
		if err != nil {
			return fmt.Errorf("could not find latest id: %w", err)
		}

		e.ID = latestID + 1
	}

	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}

	_, err := url.Parse(e.MinioLink)
	if err != nil {
		return fmt.Errorf("invalid minio link: %w", err)
	}

	_, err = url.Parse(e.YOURLSLink)
	if err != nil {
		return fmt.Errorf("invalid yourls link: %w", err)
	}

	return nil
}

var storageDir = "data"

func createStorageDirIfNotExists() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not find executable: %w", err)
	}

	storageDir = filepath.Join(filepath.Dir(exe), storageDir)

	fi, err := os.Stat(storageDir)
	if err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(storageDir, 0700)
		}

		return fmt.Errorf("could not os.Stat storage dir: %w", err)
	}

	if !fi.IsDir() {
		return errors.New("storage dir exists but is not a directory")
	}

	return nil
}

// TODO: split these into multiple files
// for performance reasons
var storageFilePath = "minls.data.json"

func createStorageFileIfNotExists() error {
	storageFilePath = filepath.Join(storageDir, storageFilePath)

	_, err := os.Stat(storageFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			_, err = os.Create(storageFilePath)
			return err
		}

		return fmt.Errorf("could not os.Stat storage file: %w", err)
	}

	return nil
}

func readCurrentDataFromFile() (*Data, error) {
	f, err := os.Open(storageFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not open storage file: %w", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("could not read from storage file: %w", err)
	}

	// if len(b) == 0 usually indicates
	// that this is the first writing to the file,
	// return default data to prevent json errors later
	if len(b) == 0 {
		return &Data{}, nil
	}

	data := &Data{}
	err = json.Unmarshal(b, data)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal data: %w", err)
	}

	return data, nil
}

var currentData *Data

const dataEntryMaxAge = 7 * 24 * time.Hour

func cleanData() error {
	if currentData == nil {
		return errors.New("current data not set up")
	}

	newData := &Data{Entries: make([]*DataEntry, 0)}
	for _, entry := range currentData.Entries {
		if time.Since(entry.Timestamp) > dataEntryMaxAge {
			continue
		}

		newData.Entries = append(newData.Entries, entry)
	}

	currentData = newData

	return nil
}

func findLatestID() (int, error) {
	if currentData == nil {
		return 0, errors.New("current data not set up")
	}

	latestID := 0
	for _, entry := range currentData.Entries {
		if entry.ID > latestID {
			latestID = entry.ID
		}
	}

	return latestID, nil
}

func writeCurrentData() error {
	if currentData == nil {
		return errors.New("current data not set up")
	}

	f, err := os.Create(storageFilePath)
	if err != nil {
		return fmt.Errorf("could not create storage file: %w", err)
	}
	defer f.Close()

	b, err := json.Marshal(currentData)
	if err != nil {
		return fmt.Errorf("could not marshal current data: %w", err)
	}

	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("could not write to storage file: %w", err)
	}

	return nil
}
