package filestorage

import (
	"bufio"
	"fmt"
	"os"
	"shorter/internal/domain"
	"shorter/internal/logger"
)

const (
	FilePerm = os.FileMode(0666)
)

type FileStorage struct {
	file    *os.File
	writer  *bufio.Writer
	scanner *bufio.Scanner
	storage domain.Storage
	logger  *logger.ZapLogger
}

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
		storage: make(domain.Storage),
		logger:  l,
	}, nil
}

func (f *FileStorage) GetInitialData() (s domain.Storage, err error) {
	storage := make(domain.Storage)
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
func (f *FileStorage) GetURL(shortURL string) (string, error) {
	return f.storage[shortURL], nil
}

func (f *FileStorage) Save(values domain.URLS) error {
	newID := values.CorrelationID
	url := values.URL
	err := f.writeToFile(newID, url)
	if err != nil {
		return err
	}
	f.storage[newID] = url
	return nil
}

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

func (f *FileStorage) SaveBatch(_ []domain.URLS) error {
	f.logger.Log.Warn("SaveBatch is not implemented")
	return fmt.Errorf("not implemented")
}

func (f *FileStorage) Ping() error {
	f.logger.Log.Warn("Ping is not implemented")
	return nil
}
