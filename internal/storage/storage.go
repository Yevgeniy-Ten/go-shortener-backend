package storage

import (
	"bufio"
	"os"
	generateRandomId "shorter/pkg"
	"sync"
)

type ShortURLStorage struct {
	storage map[string]string
	mutex   *sync.Mutex
	file    *os.File
	writer  *bufio.Writer
	scanner *bufio.Scanner
}

func NewShortURLStorage(filePath string) (*ShortURLStorage, error) {
	return &ShortURLStorage{
		storage: make(map[string]string),
		mutex:   &sync.Mutex{},
	}, nil
}

func (storage *ShortURLStorage) Save(url string) string {
	newID := generateRandomId.GenerateShortID()
	storage.mutex.Lock()
	defer storage.mutex.Unlock()
	storage.storage[newID] = url
	return newID
}

func (storage *ShortURLStorage) GetURL(id string) string {
	return storage.storage[id]
}

func (storage *ShortURLStorage) Close() error {
	return storage.file.Close()
}
