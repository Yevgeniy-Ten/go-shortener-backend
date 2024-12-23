package storage

import (
	"bufio"
	"io"
	"os"
	generateRandomId "shorter/pkg"
	"sync"
)

const (
	FilePerm = os.FileMode(0666)
)

type ShortURLStorage struct {
	storage map[string]string
	mutex   *sync.Mutex
	file    io.WriteCloser
	writer  *bufio.Writer
	scanner *bufio.Scanner
}

func NewShortURLStorage(filePath string) (*ShortURLStorage, error) {
	storage := make(map[string]string)
	if filePath == "" {
		return &ShortURLStorage{
			storage: storage,
			mutex:   &sync.Mutex{},
		}, nil
	}
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, FilePerm)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	parseLine := func(line string) (string, string) {
		return line[:8], line[9:]
	}

	for scanner.Scan() {
		line := scanner.Text()
		id, url := parseLine(line)
		storage[id] = url
	}

	return &ShortURLStorage{
		storage: storage,
		mutex:   &sync.Mutex{},
		file:    file,
		writer:  bufio.NewWriter(file),
		scanner: scanner,
	}, nil
}
func (storage *ShortURLStorage) WriteToFile(newID, url string) error {
	if storage.writer == nil {
		return nil
	}

	_, err := storage.writer.WriteString(newID + " " + url + "\n")
	if err != nil {
		return err
	}

	err = storage.writer.Flush()
	if err != nil {
		return err
	}
	return nil
}
func (storage *ShortURLStorage) Save(url string) (string, error) {
	newID := generateRandomId.GenerateShortID()
	storage.mutex.Lock()
	defer storage.mutex.Unlock()
	storage.storage[newID] = url
	err := storage.WriteToFile(newID, url)
	if err != nil {
		delete(storage.storage, newID)
		return "", err
	}
	return newID, nil
}

func (storage *ShortURLStorage) GetURL(id string) string {
	return storage.storage[id]
}

func (storage *ShortURLStorage) Close() error {
	if storage.file == nil {
		return nil
	}
	return storage.file.Close()
}
