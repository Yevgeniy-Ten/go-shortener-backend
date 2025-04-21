package filestorage

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"shorter/internal/domain"
	"shorter/internal/logger"
)

// FilePerm is a file permission
const (
	FilePerm = os.FileMode(0666)
)

// FileStorage is a struct for file storage
type FileStorage struct {
	file    *os.File
	writer  *bufio.Writer
	scanner *bufio.Scanner
	storage domain.URLStorage
	logger  *logger.ZapLogger
}

// New creates a new file storage
func New(filePath string, l *logger.ZapLogger) (*FileStorage, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path is empty")
	}
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, FilePerm)
	if err != nil {
		return nil, err
	}

	return &FileStorage{
		file:    file,
		writer:  bufio.NewWriter(file),
		scanner: bufio.NewScanner(file),
		storage: make(domain.URLStorage),
		logger:  l,
	}, nil
}

// GetInitialData returns initial data
func (f *FileStorage) GetInitialData() (s domain.URLStorage, err error) {
	storage := make(domain.URLStorage)
	parseLine := func(line string) (string, string) {
		return line[:8], line[9:]
	}

	for f.scanner.Scan() {
		line := f.scanner.Text()
		id, url := parseLine(line)
		storage[id] = url
	}
	return storage, nil
}

// GetURL returns the URL by short URL
func (f *FileStorage) GetURL(shortURL string) (string, error) {
	return f.storage[shortURL], nil
}

// Save saves the URL
func (f *FileStorage) Save(values domain.URLS, _ int) error {
	newID := values.CorrelationID
	url := values.URL
	err := f.writeToFile(newID, url)
	if err != nil {
		return err
	}
	f.storage[newID] = url
	return nil
}

// Close closes the file
func (f *FileStorage) Close() error {
	return f.file.Close()
}

func (f *FileStorage) writeToFile(newID, url string) error {
	_, err := f.writer.WriteString(newID + " " + url + "\n")
	if err != nil {
		return err
	}

	err = f.writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

// SaveBatch saves batch URLs
func (f *FileStorage) SaveBatch(_ context.Context, _ []domain.URLS, _ int) error {
	f.logger.Log.Warn("SaveBatch is not implemented")
	return errors.New("not implemented")
}

// GetUserURLs returns user URLs
func (f *FileStorage) GetUserURLs(_ int, _ string) ([]domain.UserURLs, error) {
	return nil, errors.New("not implemented")
}

// DeleteURLs deletes URLs
func (f *FileStorage) DeleteURLs(_ []string, _ int) error {
	f.logger.Log.Warn("DeleteURLs is not implemented")
	return errors.New("not implemented")
}
