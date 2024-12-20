package storage

import (
	generateRandomId "shorter/pkg"
	"sync"
)

type ShortURLStorage struct {
	storage map[string]string
	mutex   *sync.Mutex
}

func NewShortURLStorage() *ShortURLStorage {
	return &ShortURLStorage{
		storage: make(map[string]string),
		mutex:   &sync.Mutex{},
	}
}

func (storage *ShortURLStorage) Save(url string) string {
	newID := generateRandomId.GenerateShortID()
	storage.mutex.Lock()
	storage.storage[newID] = url
	defer storage.mutex.Unlock()
	return newID
}

func (storage *ShortURLStorage) GetURL(id string) string {
	return storage.storage[id]
}

var GlobalURLStorage = NewShortURLStorage()
